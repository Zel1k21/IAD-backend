package repository

import (
	"gorm.io/gorm"

	"errors"
	"iad-backend/internal/app/ds"
)

type StageRequestRepository struct {
	db *gorm.DB
}

func NewStageRequestRepository(db *gorm.DB) *StageRequestRepository {
	return &StageRequestRepository{db: db}
}

func (r *StageRequestRepository) GetStageRequestIDEntryCountByUserID(userID uint64) (int, int, error) {
	var stageRequest ds.StageRequest
	err := r.db.
		Model(&ds.StageRequest{}).
		Where("status = 1 and user_id = ?", userID).
		Take(&stageRequest).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}

	var count int64
	err = r.db.
		Model(&ds.StageRequest{}).
		Where("id = ?", stageRequest.ID).
		Joins("StageRequestToStage").
		Count(&count).Error

	if err != nil {
		return 0, 0, err
	}

	return int(count), int(stageRequest.ID), nil
}

func (r *StageRequestRepository) GetStageRequestByID(id uint64, userID uint64) (*ds.StageRequest, error) {
	var stageRequest ds.StageRequest
	err := r.db.
		Preload("StageRequestToStage").
		Preload("StageRequestToStage.Stage").
		Where("status = 1 and user_id = ?", userID).
		First(&stageRequest, id).Error

	if err != nil {
		return nil, err
	}

	return &stageRequest, nil
}

func (r *StageRequestRepository) AddStageToStageRequest(stageId uint64, userId uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var stage ds.Stage
		err := r.db.First(&stage, stageId).Error
		if err != nil {
			return err
		}

		var stageRequest ds.StageRequest
		err = r.db.
			Where("status = 1 and user_id = ?", userId).
			Take(&stageRequest).Error
		notFound := errors.Is(err, gorm.ErrRecordNotFound)
		if err != nil && !notFound {
			return err
		}

		if notFound {
			stageRequest = ds.StageRequest{User: ds.User{ID: userId}}

			err = r.db.Create(&stageRequest).Error
			if err != nil {
				return err
			}

		}

		stageRequestToStage := ds.StageRequestToStage{RequestID: stageRequest.ID, StageID: stage.ID}
		r.db.Create(&stageRequestToStage)

		return nil
	})
}

func (r *StageRequestRepository) DeleteStageRequest(requestID uint64, userID uint64) error {
	query := "update stage_request set status = 2 where status = 1 and id = $1 and user_id = $2"
	err := r.db.Exec(query, requestID, userID).Row().Err()
	if err != nil {
		return err
	}
	return nil
}
