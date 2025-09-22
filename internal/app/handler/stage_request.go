package handler

import (
	"iad-backend/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type StageRequestHandler struct {
	StageRequestRepository *repository.StageRequestRepository
	StageRepository        *repository.StageRepository
}

func NewStageRequestHandler(stageRequestRepository *repository.StageRequestRepository, stageRepository *repository.StageRepository) *StageRequestHandler {
	return &StageRequestHandler{
		StageRequestRepository: stageRequestRepository,
		StageRepository:        stageRepository,
	}
}

func (h *StageRequestHandler) GetStageRequestByID(ctx *gin.Context) {
	stageIDStr := ctx.Param("id")
	stageRequestID, err := strconv.Atoi(stageIDStr)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	stageReqView, err := h.StageRequestRepository.GetStageRequestViewByID(stageRequestID, h.StageRepository)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "stage_request.html", gin.H{
		"title": "Просмотр заявки",
		"view":  &stageReqView,
	})
}
