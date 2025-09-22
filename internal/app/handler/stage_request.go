package handler

import (
	"errors"
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

	ctx.HTML(http.StatusOK, "stage_request.html", gin.H{
		"title":        "Просмотр заявки",
		"stageRequest": &stageRequest,
	})
}

func (h *StageRequestHandler) DeleteStageRequest(ctx *gin.Context) {
	requestIDStr := ctx.PostForm("request-id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = h.repo.StageRequest.DeleteStageRequest(requestID, 1)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/stages")
}
