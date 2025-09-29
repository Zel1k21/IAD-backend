package repository

import "fmt"

type StageRequestRepository struct {
}

type CompTextInput struct {
	ShowLabel   bool
	Label       string
	Type        string
	Name        string
	Placeholder string
	Value       string
}

func NewStageRequestRepository() (*StageRequestRepository, error) {
	return &StageRequestRepository{}, nil
}

type StageRequest struct {
	ID          int
	ProductName CompTextInput
	StageResult int
}

type StageRequestViewEntry struct {
	Stage           Stage
	InputField1     int
	Field1Dimension string
	InputField2     int
	Field2Dimension string
	CardResult      int
}

type StageRequestView struct {
	StageRequest StageRequest
	Entries      []StageRequestViewEntry
}

var stageRequestViewByID = map[int]StageRequestView{
	1: {
		StageRequest: StageRequest{
			ProductName: CompTextInput{
				ShowLabel:   true,
				Type:        "text",
				Name:        "query",
				Placeholder: "Введите название этапа",
				Value:       "Парта",
			},
			StageResult: 37,
		},

		Entries: []StageRequestViewEntry{
			{
				Stage:           stages[1],
				InputField1:     3000,
				Field1Dimension: "Объём, м^3",
				InputField2:     100000,
				Field2Dimension: "Энергия, кВт·ч",
				CardResult:      20,
			},
			{
				Stage:           stages[2],
				InputField1:     23400,
				Field1Dimension: "Масса товара, т",
				InputField2:     482,
				Field2Dimension: "Расстояние, км",
				CardResult:      17,
			},
		},
	},
}

func (*StageRequestRepository) GetStageRequestEntryCountByID(id int) (int, error) {
	if len(stageRequestViewByID) == 0 {
		return 0, fmt.Errorf("array is empty")
	}

	stageRequestView, found := stageRequestViewByID[id]
	if !found {
		return 0, fmt.Errorf("stage request not found")
	}

	return len(stageRequestView.Entries), nil
}

func (*StageRequestRepository) GetStageRequestViewByID(id int, stageRepo *StageRepository) (*StageRequestView, error) {
	if len(stageRequestViewByID) == 0 {
		return nil, fmt.Errorf("array is empty")
	}

	stageRequestView, found := stageRequestViewByID[id]
	if !found {
		return nil, fmt.Errorf("stage request not found")
	}

	return &stageRequestView, nil
}
