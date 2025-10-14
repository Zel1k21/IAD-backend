package repository

import (
	"context"
	"fmt"
	"iad-backend/internal/app/dsn"
	"iad-backend/internal/app/redis"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type JWTConfig struct {
	Secret        string
	ExpiresIn     time.Duration
	RefreshIn     time.Duration
	SigningMethod jwt.SigningMethod
}

type Repository struct {
	db           *gorm.DB
	Stage        *StageRepository
	StageRequest *StageRequestRepository
	User         *UserRepository
	redis        *redis.Client
}

func (r *Repository) GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}

func (r *Repository) GetJWTConfig() *JWTConfig {
	secret := r.GetJWTSecret()
	if secret == "" {
		secret = "default-jwt-secret-key"
	}

	return &JWTConfig{
		Secret:        secret,
		ExpiresIn:     time.Hour * 24,
		RefreshIn:     time.Hour * 24 * 7,
		SigningMethod: jwt.SigningMethodHS256,
	}
}

func NewRepository() (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	minioClient, err := InitMinioClient()
	if err != nil {
		return nil, err
	}

	redisClient, err := InitRedisClient()
	if err != nil {
		return nil, err
	}

	return &Repository{
		db:           db,
		Stage:        NewStageRepository(db, minioClient),
		StageRequest: NewStageRequestRepository(db),
		User:         NewUserRepository(db),
		redis:        redisClient,
	}, nil
}

func CloseDBConn(r *Repository) {
	dbInstance, _ := r.db.DB()
	_ = dbInstance.Close()

	if r.redis != nil {
		_ = r.redis.Close()
	}
}

func (r *Repository) GetRedisClient() *redis.Client {
	return r.redis
}

func InitRedisClient() (*redis.Client, error) {
	host := os.Getenv("REDIS_HOST")
	portStr := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	user := os.Getenv("REDIS_USER")

	if host == "" {
		host = "localhost"
	}

	port := 6379
	if portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("redis port must be int value: %v", err)
		}
	}

	cfg := &redis.Config{
		Host:        host,
		Port:        port,
		Password:    password,
		User:        user,
		DialTimeout: 10 * time.Second,
		ReadTimeout: 30 * time.Second,
	}

	client := redis.New(cfg)

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return client, nil
}

func InitMinioClient() (*minio.Client, error) {
	endpoint := os.Getenv("MINIO_HOST") + ":" + os.Getenv("MINIO_SERVER_PORT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %v", err)
	}

	ctx := context.Background()

	exists, err := minioClient.BucketExists(ctx, stageImagesBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %v", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, stageImagesBucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %v", err)
		}
		logrus.Printf("Bucket '%s' created successfully\n", stageImagesBucket)
	}

	return minioClient, nil
}
