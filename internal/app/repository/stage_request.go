package repository

import (
	"time"

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

func (r *StageRequestRepository) GetDraftRequestInfo(userID uint64) (uint64, int, error) {
	var request ds.StageRequest
	err := r.db.
		Where("status = 1 AND user_id = ?", userID).
		First(&request).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}

	var count int64
	err = r.db.
		Model(&ds.StageRequestToStage{}).
		Where("request_id = ?", request.ID).
		Count(&count).Error

	if err != nil {
		return 0, 0, err
	}

	return request.ID, int(count), nil
}

func (r *StageRequestRepository) GetStageRequests(statusFilter uint8, dateFrom, dateTo *time.Time) ([]ds.StageRequest, error) {
	var requests []ds.StageRequest

	query := r.db.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username")
		}).
		Preload("Morderator", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username")
		}).
		Where("status != 1 AND status != 2") // исключаем черновики и удалённые

	if statusFilter != 0 {
		query = query.Where("status = ?", statusFilter)
	}

	if dateFrom != nil {
		query = query.Where("formed_at >= ?", dateFrom)
	}
	if dateTo != nil {
		query = query.Where("formed_at <= ?", dateTo)
	}

	err := query.Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *StageRequestRepository) GetStageRequestByID(id uint64, userID uint64) (*ds.StageRequest, error) {
	var request ds.StageRequest
	err := r.db.
		Preload("StageRequestToStage").
		Preload("StageRequestToStage.Stage").
		Where("status != 2 AND user_id = ?", userID).
		First(&request, id).Error

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *StageRequestRepository) UpdateStageRequest(id uint64, productName *string) error {
	updates := make(map[string]interface{})

	if productName != nil {
		updates["product_name"] = *productName
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.
		Model(&ds.StageRequest{}).
		Where("id = ? AND status != 2", id).
		Updates(updates).Error
}

func (r *StageRequestRepository) FormRequest(id uint64, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var request ds.StageRequest
		err := tx.
			Preload("StageRequestToStage").
			Where("id = ? AND user_id = ? AND status = 1", id, userID).
			First(&request).Error

		if err != nil {
			return err
		}

		if request.ProductName == "" {
			return errors.New("product name is required")
		}
		if len(request.StageRequestToStage) == 0 {
			return errors.New("at least one stage is required")
		}

		return tx.Model(&request).Updates(map[string]interface{}{
			"status":    3,
			"formed_at": time.Now(),
		}).Error
	})
}

func (r *StageRequestRepository) ResolveOrRejectRequest(id uint64, moderatorID uint64, status uint8) error {
	if status != 4 && status != 5 {
		return errors.New("invalid status for moderator action")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		var request ds.StageRequest
		err := tx.
			Preload("StageRequestToStage").
			Preload("StageRequestToStage.Stage").
			Where("id = ? AND status = 3", id).
			First(&request).Error

		if err != nil {
			return err
		}

		for i := range request.StageRequestToStage {
			stageRequestToStage := &request.StageRequestToStage[i]
			stage := stageRequestToStage.Stage
			stageEmission := uint64(
				float64(stageRequestToStage.InputField1)*stage.FirstDimensionConst +
					float64(stageRequestToStage.InputField2)*stage.SecondDimensionConst,
			)

			err := tx.Model(stageRequestToStage).
				Update("stage_calculation_result", stageEmission).Error
			if err != nil {
				return err
			}
		}

		calculatedEmission := r.CalculateEmission(request.ID)

		updates := map[string]interface{}{
			"status":                      status,
			"moderator_id":                moderatorID,
			"closed_at":                   time.Now(),
			"emission_calculation_result": calculatedEmission,
		}

		return tx.Model(&request).Updates(updates).Error
	})
}

func (r *StageRequestRepository) DeleteRequest(id uint64, userID uint64) error {
	return r.db.
		Model(&ds.StageRequest{}).
		Where("id = ? AND user_id = ? AND status = 1", id, userID).
		Update("status", 2).Error
}

func (r *StageRequestRepository) RemoveStageFromRequest(requestID uint64, stageID uint64, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var request ds.StageRequest
		err := tx.
			Where("id = ? AND user_id = ? AND status = 1", requestID, userID).
			First(&request).Error

		if err != nil {
			return err
		}

		return tx.
			Where("request_id = ? AND stage_id = ?", requestID, stageID).
			Delete(&ds.StageRequestToStage{}).Error
	})
}

func (r *StageRequestRepository) UpdateRequestToStage(requestID uint64, stageID uint64, userID uint64, input1, input2 *uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var request ds.StageRequest
		err := tx.
			Where("id = ? AND user_id = ? AND status = 1", requestID, userID).
			First(&request).Error

		if err != nil {
			return err
		}

		updates := make(map[string]interface{})
		if input1 != nil {
			updates["input_field1"] = *input1
		}
		if input2 != nil {
			updates["input_field2"] = *input2
		}

		if len(updates) == 0 {
			return nil
		}

		return tx.
			Model(&ds.StageRequestToStage{}).
			Where("request_id = ? AND stage_id = ?", requestID, stageID).
			Updates(updates).Error
	})
}

func (r *StageRequestRepository) CalculateEmission(requestID uint64) uint64 {
	var stageRequestToStages []ds.StageRequestToStage

	err := r.db.
		Preload("Stage").
		Where("request_id = ?", requestID).
		Find(&stageRequestToStages).Error

	if err != nil {
		return 0
	}

	var totalEmission uint64 = 0

	for _, stageRequestToStage := range stageRequestToStages {
		stage := stageRequestToStage.Stage
		stageEmission := uint64(
			float64(stageRequestToStage.InputField1)*stage.FirstDimensionConst +
				float64(stageRequestToStage.InputField2)*stage.SecondDimensionConst,
		)
		totalEmission += stageEmission
	}

	return totalEmission
}

func (r *StageRequestRepository) AddStageToStageRequest(stageID uint64, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var stage ds.Stage
		err := tx.First(&stage, stageID).Error
		if err != nil {
			return err
		}

		var stageRequest ds.StageRequest
		err = tx.
			Where("status = 1 AND user_id = ?", userID).
			Take(&stageRequest).Error

		notFound := errors.Is(err, gorm.ErrRecordNotFound)
		if err != nil && !notFound {
			return err
		}

		if notFound {
			stageRequest = ds.StageRequest{
				UserID: userID,
			}
			err := tx.Create(&stageRequest).Error
			if err != nil {
				return err
			}
		}

		stageRequestToStage := ds.StageRequestToStage{
			RequestID: stageRequest.ID,
			StageID:   stageID,
		}
		return tx.Create(&stageRequestToStage).Error
	})
}
