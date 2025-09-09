package handler

import (
	"iad-backend/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type StageHandler struct {
	StageRepositiry *repository.StageRepository
}

func NewStageHandler(stageRepository *repository.StageRepository) *StageHandler {
	return &StageHandler{
		StageRepositiry: stageRepository,
	}
}

func (h *StageHandler) GetStageByID(ctx *gin.Context) {
	stageIDStr := ctx.Param("id")
	stageID, err := strconv.Atoi(stageIDStr)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	stage, err := h.StageRepositiry.GetStageByID(stageID)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "stage.html", gin.H{
		"title": "Этап жизненного цикла",
		"stage": stage,
	})
}
