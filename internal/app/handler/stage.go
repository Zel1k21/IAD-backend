package handler

import (
	"iad-backend/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type StageHandler struct {
	repo *repository.Repository
}

func NewStageHandler(repository *repository.Repository) *StageHandler {
	return &StageHandler{repo: repository}
}

func (h *StageHandler) Register(router *gin.Engine) {
	router.GET("/stage/:id", h.GetStageByID)
}

func (h *StageHandler) GetStageByID(ctx *gin.Context) {
	stageIDStr := ctx.Param("id")
	stageID, err := strconv.ParseUint(stageIDStr, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	stage, err := h.repo.Stage.GetStageByID(stageID)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "stage.html", gin.H{
		"title": "Этап жизненного цикла",
		"stage": *stage,
	})
}
