package main

import (
	"iad-backend/internal/app/ds"
	"iad-backend/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load("deploy/.env")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(
		&ds.Stage{},
		&ds.StageRequest{},
		&ds.StageRequestToStage{},
		&ds.User{},
	)
	if err != nil {
		panic(err)
	}
}
