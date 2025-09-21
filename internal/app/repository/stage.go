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
		Description: "Преобразование добытых природных ресурсов в материалы, пригодные для промышленного производства. На этом этапе сырье (руда, нефть, древесина, хлопок) очищается, сортируется и видоизменяется для создания полуфабрикатов.",
	},
	{
		ID:          3,
		Title:       "Производство готовой продукции",
		ImageURL:    "http://localhost:9000/stageimages/production_stage.png",
		Description: "Процесс создания конечного товара из переработанного сырья и полуфабрикатов. Включает проектирование, сборку, тестирование и упаковку. Здесь материалы обретают свою окончательную форму и потребительские свойства.",
	},
	{
		ID:          4,
		Title:       "Транспортировка",
		ImageURL:    "http://localhost:9000/stageimages/transportation_stage.png",
		Description: "Логистическое перемещение товаров от производителя к точкам распределения и продажи (склады, магазины). Обеспечивает доступность продукции для потребителя и включает все виды перевозок.",
	},
	{
		ID:          5,
		Title:       "Хранение и обслуживание ",
		ImageURL:    "http://localhost:9000/stageimages/storage_stage.png",
		Description: "Содержание товаров на складах и в точках продаж в условиях, обеспечивающих их сохранность и качество. Включает складскую логистику, управление запасами, маркировку, а также предпродажную подготовку и обслуживание.",
	},
	{
		ID:          6,
		Title:       "Утилизация",
		ImageURL:    "http://localhost:9000/stageimages/recycling_stage.png",
		Description: "Завершающий этап жизненного цикла, на котором товар, исчерпавший свой ресурс, перерабатывается для повторного использования материалов или безопасно уничтожается.",
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
