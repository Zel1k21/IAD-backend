package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

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

func (r *StageRequestRepository) GetStageRequests(userID uint64, isModerator bool, statusFilter uint8, dateFrom, dateTo string) ([]ds.StageRequest, []string, error) {
	var requests []ds.StageRequest
	var usernames []string

	query := r.db.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username")
		}).
		Preload("Moderator", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username")
		}).
		Where("status > 2").Order("status ASC")

	if !isModerator {
		query = query.Where("user_id = ?", userID)
	}

	if statusFilter != 0 {
		query = query.Where("status = ?", statusFilter)
	}

	if dateFrom != "" {
		query = query.Where("created_at >= ? or formed_at >= ? or closed_at >= ?", dateFrom, dateFrom, dateFrom)
	}
	if dateTo != "" {
		query = query.Where("created_at <= ? or formed_at <= ? or closed_at <= ?", dateTo, dateTo, dateTo)
	}

	err := query.Find(&requests).Error
	if err != nil {
		return nil, nil, err
	}

	for _, request := range requests {
		usernames = append(usernames, request.User.Username)
	}

	return requests, usernames, nil
}

func (r *StageRequestRepository) GetStageRequestByID(id uint64, userID uint64, isModerator bool) (*ds.StageRequest, error) {
	var request ds.StageRequest
	query := r.db.
		Preload("StageRequestToStage").
		Preload("StageRequestToStage.Stage").
		Where("status != 2")

	if !isModerator {
		query = query.Where("user_id = ?", userID)
	}

	err := query.First(&request, id).Error

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *StageRequestRepository) UpdateStageRequest(id uint64, productName *string) error {
	updates := make(map[string]any)

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

		return tx.Model(&request).Updates(map[string]any{
			"status":    3,
			"formed_at": time.Now(),
		}).Error
	})
}

func (r *StageRequestRepository) ResolveOrRejectRequest(id uint64, moderatorID uint64, status uint8) error {
	if status != 4 && status != 5 {
		return errors.New("invalid status for moderator action")
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var request ds.StageRequest
		err := tx.
			Preload("StageRequestToStage").
			Preload("StageRequestToStage.Stage").
			Where("id = ? AND status = 3", id).
			First(&request).Error

		if err != nil {
			return err
		}
		logrus.Info(request.ID)

		err = r.db.Model(&ds.StageRequest{}).
			Where("id = ?", request.ID).
			Update("status", status).Error

		return r.SendDataToCalcService(request.ID, os.Getenv("ASYNC_TOKEN"))
	})

	if err != nil {
		return err
	}

	return nil
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

		updates := make(map[string]any)
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

type CalcRequest struct {
	RequestID   uint64    `json:"id"`
	AuthToken   string    `json:"auth_token"`
	InputFields []float64 `json:"input_fields"`
	Constants   []float64 `json:"stage_constants"`
}

func (r *StageRequestRepository) SendDataToCalcService(requestID uint64, authToken string) error {
	var stageRequestToStages []ds.StageRequestToStage
	var calcServiceURL = os.Getenv("CALC_SERVICE_URL")
	err := r.db.
		Preload("Stage").
		Where("request_id = ?", requestID).
		Find(&stageRequestToStages).Error

	if err != nil {
		return err
	}

	var inputFields []float64
	var constants []float64

	for _, s := range stageRequestToStages {
		inputFields = append(inputFields, s.InputField1, s.InputField2)
		constants = append(constants, s.Stage.FirstDimensionConst, s.Stage.SecondDimensionConst)
	}

	payload := CalcRequest{
		RequestID:   requestID,
		AuthToken:   authToken,
		InputFields: inputFields,
		Constants:   constants,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(calcServiceURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode <= 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (r *StageRequestRepository) UpdateEmissionCalculationResult(requestId uint64, result float64) error {
	return r.db.
		Model(&ds.StageRequest{}).
		Where("id = ?", requestId).
		Update("calculation_result", result).Error
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
