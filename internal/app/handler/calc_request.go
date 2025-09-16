package handler

import (
	"iad-backend/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CalcRequestHandler struct {
	CalcRequestRepository *repository.CalcRequestRepository
	StageRepository       *repository.StageRepository
}

func NewCalcRequestHandler(calcRequestRepository *repository.CalcRequestRepository, stageRepository *repository.StageRepository) *CalcRequestHandler {
	return &CalcRequestHandler{
		CalcRequestRepository: calcRequestRepository,
		StageRepository:       stageRepository,
	}
}

func (h *CalcRequestHandler) GetCalcRequestByID(ctx *gin.Context) {
	stageIDStr := ctx.Param("id")
	calcRequestID, err := strconv.Atoi(stageIDStr)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	calcReqView, err := h.CalcRequestRepository.GetCalcRequestViewByID(calcRequestID, h.StageRepository)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "calc_request.html", gin.H{
		"title": "Просмотр заявки",
		"view":  &calcReqView,
	})
}
