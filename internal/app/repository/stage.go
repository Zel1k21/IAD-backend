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
		ImageURL:    "http://localhost:9000/stageimages/extraction_stage.png",
		Description: "Добыча природных ресурсов: металлов, нефти, древесины, сельхозсырья и т.д. Первичная переработка (очистка, сортировка, подготовка к производству).",
	},
	{
		ID:          2,
		Title:       "Переработка сырья",
		ImageURL:    "http://localhost:9000/stageimages/transformation_stage.png",
		Description: "Stage 2 description",
	},
	{
		ID:          3,
		Title:       "Производство готовой продукции",
		ImageURL:    "http://localhost:9000/stageimages/production_stage.png",
		Description: "Stage 3 description",
	},
	{
		ID:          4,
		Title:       "Транспортировка",
		ImageURL:    "http://localhost:9000/stageimages/transportation_stage.png",
		Description: "Stage 4 description",
	},
	{
		ID:          5,
		Title:       "Хранение и обслуживание ",
		ImageURL:    "http://localhost:9000/stageimages/storage_stage.png",
		Description: "Stage 5 description",
	},
	{
		ID:          6,
		Title:       "Утилизация",
		ImageURL:    "http://localhost:9000/stageimages/recycling_stage.png",
		Description: "Stage 6 description",
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
