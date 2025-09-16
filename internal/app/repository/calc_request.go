package repository

import "fmt"

type CalcRequestRepository struct {
}

type CompTextInput struct {
	ShowLabel   bool
	Label       string
	Type        string
	Name        string
	Placeholder string
	Value       string
}

func NewCalcRequestRepository() (*CalcRequestRepository, error) {
	return &CalcRequestRepository{}, nil
}

type CalcRequest struct {
	ID                int
	ProductName       CompTextInput
	CalculationResult int
}

type CalcRequestToStage struct {
	RequestID       int
	StageID         int
	InputField1     int
	Field1Dimension string
	InputField2     int
	Field2Dimension string
	CardResult      int
}

type CalcRequestViewEntry struct {
	Stage           Stage
	InputField1     int
	Field1Dimension string
	InputField2     int
	Field2Dimension string
	CardResult      int
}

type CalcRequestView struct {
	CalcRequest CalcRequest
	Entries     []CalcRequestViewEntry
}

var calcRequests = []CalcRequest{
	{
		ID: 1,
		ProductName: CompTextInput{
			ShowLabel:   true,
			Type:        "text",
			Name:        "query",
			Placeholder: "Введите название этапа",
			Value:       "Парта",
		},
		CalculationResult: 37,
	},
}

var CalcRequestToStages = []CalcRequestToStage{
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

func (*CalcRequestRepository) GetCalcRequestEntryCountByID(id int) (int, error) {
	if len(calcRequests) == 0 {
		return 0, fmt.Errorf("array is empty")
	}

	var calcRequest *CalcRequest = nil
	for _, req := range calcRequests {
		if req.ID == id {
			calcRequest = &req
		}
	}

	if calcRequest == nil {
		return 0, fmt.Errorf("calc request not found")
	}

	var calcRequestEntryCount int = 0
	for _, reqToStage := range CalcRequestToStages {
		if reqToStage.RequestID == id {
			calcRequestEntryCount++
		}
	}

	return calcRequestEntryCount, nil
}

func (*CalcRequestRepository) GetCalcRequestViewByID(id int, stageRepo *StageRepository) (*CalcRequestView, error) {
	if len(calcRequests) == 0 {
		return nil, fmt.Errorf("array is empty")
	}

	var calcRequest *CalcRequest = nil
	for _, req := range calcRequests {
		if req.ID == id {
			calcRequest = &req
		}
	}

	if calcRequest == nil {
		return nil, fmt.Errorf("calc request not found")
	}

	calcRequestView := CalcRequestView{
		CalcRequest: *calcRequest,
	}

	for _, reqToStage := range CalcRequestToStages {
		if reqToStage.RequestID == id {
			stage, err := stageRepo.GetStageByID(reqToStage.StageID)
			if err != nil {
				return nil, err
			}
			calcRequestView.Entries = append(calcRequestView.Entries, CalcRequestViewEntry{
				Stage:           *stage,
				InputField1:     reqToStage.InputField1,
				InputField2:     reqToStage.InputField2,
				Field1Dimension: reqToStage.Field1Dimension,
				Field2Dimension: reqToStage.Field2Dimension,
				CardResult:      reqToStage.CardResult,
			})
		}
	}

	return &calcRequestView, nil
}
