package handler

import (
	"iad-backend/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RequestStageHandler struct {
	repo *repository.Repository
}

func NewStageRequestToStageHandler(repo *repository.Repository) *RequestStageHandler {
	return &RequestStageHandler{
		repo: repo,
	}
}

type UpdateStageToRequestConnection struct {
	InputField1 *uint64 `json:"input_field_1"`
	InputField2 *uint64 `json:"input_field_2"`
}

func (h *RequestStageHandler) RemoveStageToRequestConnection(ctx *gin.Context) {
	requestIDStr := ctx.Param("requestId")
	stageIDStr := ctx.Param("stageId")

	requestID, err := strconv.ParseUint(requestIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	stageID, err := strconv.ParseUint(stageIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage ID"})
		return
	}

	userID := GetFixedUserID()
	if err := h.repo.StageRequest.RemoveStageFromRequest(requestID, stageID, userID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove stage from request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Stage removed from request successfully"})
}

func (h *RequestStageHandler) UpdateStageToRequestConnection(ctx *gin.Context) {
	requestIDStr := ctx.Param("requestId")
	stageIDStr := ctx.Param("stageId")

	requestID, err := strconv.ParseUint(requestIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	stageID, err := strconv.ParseUint(stageIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage ID"})
		return
	}

	var req UpdateStageToRequestConnection
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage request data"})
		return
	}

	userID := GetFixedUserID()
	if err := h.repo.StageRequest.UpdateRequestToStage(requestID, stageID, userID, req.InputField1, req.InputField2); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request stage"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request stage updated successfully"})
}
