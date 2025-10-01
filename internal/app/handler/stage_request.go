package handler

import (
	"errors"
	"iad-backend/internal/app/ds"
	"iad-backend/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StageRequestHandler struct {
	repo *repository.Repository
}

func NewStageRequestHandler(repository *repository.Repository) *StageRequestHandler {
	return &StageRequestHandler{repo: repository}
}

func (h *StageRequestHandler) Register(router *gin.Engine) {
	router.GET("/stage_request/:id", h.GetStageRequestByID)
	router.POST("/stage_request/:id", h.DeleteStageRequest)
}

type StageRequestTemplateEntry struct {
	Stage           ds.Stage
	Field1Dimension string
	Field2Dimension string
	InputField1     CompTextInput
	InputField2     CompTextInput
	CardResult      uint64
}

func NewStageRequestTemplateEntry(stageReqEntry *ds.StageRequestToStage) *StageRequestTemplateEntry {
	return &StageRequestTemplateEntry{
		Stage:           stageReqEntry.Stage,
		Field1Dimension: stageReqEntry.Stage.FirstDimensionName,
		Field2Dimension: stageReqEntry.Stage.SecondDimensionName,
		InputField1: CompTextInput{
			Value:       strconv.FormatUint(stageReqEntry.InputField1, 10),
			Placeholder: "Введите значение",
		},
		InputField2: CompTextInput{
			Value:       strconv.FormatUint(stageReqEntry.InputField2, 10),
			Placeholder: "Введите значение",
		},
		CardResult: stageReqEntry.StageCalculationResult,
	}
}

type StageRequestTemplate struct {
	ID                        uint64
	ProductName               CompTextInput
	Entries                   []StageRequestTemplateEntry
	EmissionCalculationResult uint64
}

func (h *StageRequestHandler) GetStageRequestByID(ctx *gin.Context) {
	stageIDStr := ctx.Param("id")
	reqID, err := strconv.ParseUint(stageIDStr, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	stageRequest, err := h.repo.StageRequest.GetStageRequestByID(reqID, 1)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	stageReqTemplate := StageRequestTemplate{
		ID: stageRequest.ID,
		ProductName: CompTextInput{
			Value:       stageRequest.ProductName,
			Placeholder: "Название продукта",
		},
		EmissionCalculationResult: stageRequest.EmissionCalculationResult,
	}

	for _, stageReqToLamp := range stageRequest.StageRequestToStage {
		stageReqTemplate.Entries = append(stageReqTemplate.Entries, *NewStageRequestTemplateEntry(&stageReqToLamp))
	}

	ctx.HTML(http.StatusOK, "stage_request.html", gin.H{
		"title":        "Просмотр заявки",
		"stageRequest": &stageReqTemplate,
	})
}

func (h *StageRequestHandler) DeleteStageRequest(ctx *gin.Context) {
	requestIDStr := ctx.PostForm("request-id")
	logrus.Infof("DeleteStageRequest called with request-id: %s", requestIDStr)

	if requestIDStr == "" {
		logrus.Error("Empty request-id parameter")
		ctx.Status(http.StatusBadRequest)
		return
	}

	requestID, err := strconv.ParseUint(requestIDStr, 10, 64)
	if err != nil {
		logrus.Errorf("Failed to parse request ID: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	logrus.Infof("Attempting to delete stage request ID: %d for user ID: 1", requestID)
	err = h.repo.StageRequest.DeleteStageRequest(requestID, 1)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Errorf("Stage request not found: %v", err)
			ctx.Status(http.StatusNotFound)
			return
		}
		logrus.Errorf("Failed to delete stage request: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	logrus.Infof("Successfully deleted stage request ID: %d, redirecting to /stages", requestID)
	ctx.Redirect(http.StatusSeeOther, "/stages")
}
