package repository

import (
	"gorm.io/gorm"

	"iad-backend/internal/app/ds"
)

type StageRepository struct {
	db *gorm.DB
}

func NewStageRepository(db *gorm.DB) *StageRepository {
	return &StageRepository{db: db}
}

func (r *StageRepository) GetStageByID(id uint64) (*ds.Stage, error) {
	var stage ds.Stage
	err := r.db.
		Where("is_deleted = false").
		First(&stage, id).Error

	if err != nil {
		return nil, err
	}
	return &stage, err
}

func (r *StageRepository) GetStages() ([]ds.Stage, error) {
	var stages []ds.Stage
	err := r.db.
		Where("is_deleted = false").
		Find(&stages).Error

	if err != nil {
		return nil, err
	}
	return stages, err
}

func (r *StageRepository) GetStagesByTitle(title string) ([]ds.Stage, error) {
	var stages []ds.Stage
	err := r.db.
		Where("is_deleted = false and title like ?", "%"+title+"%").
		Find(&stages).Error

	if err != nil {
		return nil, err
	}
	return stages, err
}
