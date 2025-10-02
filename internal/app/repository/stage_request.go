package repository

import (
	"github.com/sirupsen/logrus"
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

	return int(stageRequest.ID), int(count), nil
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
		err := tx.First(&stage, stageId).Error
		if err != nil {
			return err
		}

		var stageRequest ds.StageRequest
		err = tx.
			Where("status = 1 and user_id = ?", userId).
			Take(&stageRequest).Error
		notFound := errors.Is(err, gorm.ErrRecordNotFound)
		if err != nil && !notFound {
			return err
		}

		if notFound {
			stageRequest = ds.StageRequest{User: ds.User{ID: userId}}

			err = tx.Create(&stageRequest).Error
			if err != nil {
				return err
			}

		}

		var existingStageRequestToStage ds.StageRequestToStage
		err = tx.
			Where("request_id = ? AND stage_id = ?", stageRequest.ID, stageId).
			First(&existingStageRequestToStage).Error

		if err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		stageRequestToStage := ds.StageRequestToStage{RequestID: stageRequest.ID, StageID: stageId}
		err = tx.Create(&stageRequestToStage).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *StageRequestRepository) DeleteStageRequest(requestID uint64, userID uint64) error {
	logrus.Infof("Executing SQL UPDATE for stage request ID: %d, user ID: %d", requestID, userID)

	query := "UPDATE stage_requests SET status = 2 WHERE status = 1 AND id = ? AND user_id = ?"
	result := r.db.Exec(query, requestID, userID)

	if result.Error != nil {
		logrus.Errorf("SQL UPDATE failed: %v", result.Error)
		return result.Error
	}

	rowsAffected := result.RowsAffected
	logrus.Infof("SQL UPDATE completed, rows affected: %d", rowsAffected)

	if rowsAffected == 0 {
		logrus.Warnf("No rows affected - stage request not found or already deleted (ID: %d, User: %d)", requestID, userID)
		return gorm.ErrRecordNotFound
	}

	logrus.Infof("Successfully updated stage request status to 2 for ID: %d", requestID)
	return nil
}
