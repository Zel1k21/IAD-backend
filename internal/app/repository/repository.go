package repository

import (
	"iad-backend/internal/app/dsn"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	Stage        *StageRepository
	StageRequest *StageRequestRepository
	User         *UserRepository
}

func NewRepository() (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		Stage:        NewStageRepository(db),
		StageRequest: NewStageRequestRepository(db),
		User:         NewUserRepository(db),
	}, nil
}
