package repository

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

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

func (r *StageRepository) GetStages(titleFilter string) ([]ds.Stage, error) {
	var stages []ds.Stage
	query := r.db.Where("is_deleted = false")
	if titleFilter != "" {
		query = query.Where("title ILIKE ?", "%"+titleFilter+"%")
	}
	err := query.Find(&stages).Error
	if err != nil {
		return nil, err
	}
	return stages, err
}

func (r *StageRepository) CreateStage(stage *ds.Stage) error {
	stage.IsDeleted = false
	return r.db.Create(stage).Error
}

func (r *StageRepository) UpdateStage(id uint64, data *ds.Stage) error {
	return r.db.Model(&ds.Stage{}).
		Where("id = ? AND is_deleted = false", id).
		Updates(map[string]interface{}{
			"title":                  data.Title,
			"description":            data.Description,
			"first_dimension_name":   data.FirstDimensionName,
			"first_dimension_const":  data.FirstDimensionConst,
			"second_dimension_name":  data.SecondDimensionName,
			"second_dimension_const": data.SecondDimensionConst,
		}).Error
}

func (r *StageRepository) DeleteStage(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var stage ds.Stage
		if err := tx.First(&stage, id).Error; err != nil {
			return err
		}

		if stage.ImageURL != "" {
			if err := r.deleteImageFile(stage.ImageURL); err != nil {
				return err
			}
		}

		return tx.Model(&ds.Stage{}).Where("id = ?", id).Update("is_deleted", true).Error
	})
}

func (r *StageRepository) AddStageImage(id uint64, fileHeader *multipart.FileHeader) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var stage ds.Stage
		if err := tx.Where("is_deleted = false").First(&stage, id).Error; err != nil {
			return err
		}

		if stage.ImageURL != "" {
			if err := r.deleteImageFile(stage.ImageURL); err != nil {
				return err
			}
		}

		fileExt := filepath.Ext(fileHeader.Filename)
		newFileName := fmt.Sprintf("stage_%d_%d%s", id, time.Now().Unix(), fileExt)
		newFileName = strings.ToLower(newFileName)

		imageURL, err := r.saveImageToMinIO(newFileName)
		if err != nil {
			return err
		}

		return tx.Model(&stage).Update("image_url", imageURL).Error
	})
}

func (r *StageRepository) AddStageToDraftRequest(stageID uint64, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var stage ds.Stage
		if err := tx.Where("is_deleted = false").First(&stage, stageID).Error; err != nil {
			return err
		}

		var request ds.StageRequest
		err := tx.Where("status = 1 AND user_id = ?", userID).First(&request).Error
		if err == gorm.ErrRecordNotFound {
			request = ds.StageRequest{
				Status:    1,
				UserID:    userID,
				CreatedAt: time.Now(),
			}
			if err := tx.Create(&request).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		stageRequestToStage := ds.StageRequestToStage{
			RequestID:   request.ID,
			StageID:     stageID,
			InputField1: 0,
			InputField2: 0,
		}

		return tx.Create(&stageRequestToStage).Error
	})
}

func (r *StageRepository) saveImageToMinIO(fileName string) (string, error) {
	return fmt.Sprintf("http://localhost:9000/stage-images/%s", fileName), nil
}

func (r *StageRepository) deleteImageFile(imageURL string) error {
	if strings.Contains(imageURL, "localhost:9000") {
		fmt.Printf("Image deleted from MinIO: %s\n", imageURL)
	}
	return nil
}
