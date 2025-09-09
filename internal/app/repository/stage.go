package repository

import (
	"fmt"
	"strings"
)

type StageRepository struct {
}

func NewStageRepository() (*StageRepository, error) {
	return &StageRepository{}, nil
}

type Stage struct {
	ID          int
	Title       string
	ImageURL    string
	Description string
}

var stages = []Stage{
	{
		ID:          1,
		Title:       "Добыча и подготовка сырья ",
		ImageURL:    "/static/stages/1.jpg",
		Description: "Stage 1 description",
	},
	{
		ID:          2,
		Title:       "Переработка сырья",
		ImageURL:    "/static/stages/2.jpg",
		Description: "Stage 2 description",
	},
	{
		ID:          3,
		Title:       "Производство готовой продукции",
		ImageURL:    "/static/stages/3.jpg",
		Description: "Stage 3 description",
	},
}

func (r *StageRepository) GetStageByID(id int) (*Stage, error) {
	if len(stages) == 0 {
		return nil, fmt.Errorf("array is empty")
	}
	for _, stage := range stages {
		if stage.ID == id {
			return &stage, nil
		}
	}
	return nil, fmt.Errorf("stage not found")
}

func (r *StageRepository) GetStages() ([]Stage, error) {
	if len(stages) == 0 {
		return nil, fmt.Errorf("array is empty")
	}
	return stages, nil
}

func (r *StageRepository) GetStagesByTitle(title string) ([]Stage, error) {
	stages, err := r.GetStages()
	if err != nil {
		return []Stage{}, err
	}

	var result []Stage
	for _, stage := range stages {
		if strings.Contains(strings.ToLower(stage.Title), strings.ToLower(title)) {
			result = append(result, stage)
		}
	}
	return result, nil
}
