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
	ID                 int
	ProductName        CompTextInput
	StageulationResult int
}

type StageRequestToStage struct {
	RequestID       int
	StageID         int
	InputField1     int
	Field1Dimension string
	InputField2     int
	Field2Dimension string
	CardResult      int
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

var stageRequests = []StageRequest{
	{
		ID: 1,
		ProductName: CompTextInput{
			ShowLabel:   true,
			Type:        "text",
			Name:        "query",
			Placeholder: "Введите название этапа",
			Value:       "Парта",
		},
		StageulationResult: 37,
	},
}

var StageRequestToStages = []StageRequestToStage{
	{
		RequestID:       1,
		StageID:         1,
		InputField1:     3000,
		Field1Dimension: "Объём, м^3",
		InputField2:     100000,
		Field2Dimension: "Энергия, кВт·ч",
		CardResult:      20,
	},
	{
		RequestID:       1,
		StageID:         2,
		InputField1:     23400,
		Field1Dimension: "Масса товара, т",
		InputField2:     482,
		Field2Dimension: "Расстояние, км",
		CardResult:      17,
	},
}

func (*StageRequestRepository) GetStageRequestEntryCountByID(id int) (int, error) {
	if len(stageRequests) == 0 {
		return 0, fmt.Errorf("array is empty")
	}

	var stageRequest *StageRequest = nil
	for _, req := range stageRequests {
		if req.ID == id {
			stageRequest = &req
		}
	}

	if stageRequest == nil {
		return 0, fmt.Errorf("stage request not found")
	}

	var stageRequestEntryCount int = 0
	for _, reqToStage := range StageRequestToStages {
		if reqToStage.RequestID == id {
			stageRequestEntryCount++
		}
	}

	return stageRequestEntryCount, nil
}

func (*StageRequestRepository) GetStageRequestViewByID(id int, stageRepo *StageRepository) (*StageRequestView, error) {
	if len(stageRequests) == 0 {
		return nil, fmt.Errorf("array is empty")
	}

	var stageRequest *StageRequest = nil
	for _, req := range stageRequests {
		if req.ID == id {
			stageRequest = &req
		}
	}

	if stageRequest == nil {
		return nil, fmt.Errorf("stage request not found")
	}

	stageRequestView := StageRequestView{
		StageRequest: *stageRequest,
	}

	for _, reqToStage := range StageRequestToStages {
		if reqToStage.RequestID == id {
			stage, err := stageRepo.GetStageByID(reqToStage.StageID)
			if err != nil {
				return nil, err
			}
			stageRequestView.Entries = append(stageRequestView.Entries, StageRequestViewEntry{
				Stage:           *stage,
				InputField1:     reqToStage.InputField1,
				InputField2:     reqToStage.InputField2,
				Field1Dimension: reqToStage.Field1Dimension,
				Field2Dimension: reqToStage.Field2Dimension,
				CardResult:      reqToStage.CardResult,
			})
		}
	}

	return &stageRequestView, nil
}
