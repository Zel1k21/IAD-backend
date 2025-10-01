package handler

import (
	"iad-backend/internal/app/repository"
	"net/http"

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

type RemoveStageToRequestConnection struct {
	StageRequestID uint64 `json:"stage_request_id" binding:"required"`
	StageID        uint64 `json:"stage_id" binding:"required"`
}

type UpdateStageToRequestConnection struct {
	StageRequestID uint64 `json:"stage_request_id" binding:"required"`
	StageID        uint64 `json:"stage_id" binding:"required"`
	InputField1    uint64 `json:"input_field_1"`
	InputField2    uint64 `json:"input_field_2"`
}

func (h *RequestStageHandler) RemoveStageToRequestConnection(ctx *gin.Context) {
	var req RemoveStageToRequestConnection
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage request data"})
		return
	}

	userID := GetFixedUserID()
	if err := h.repo.StageRequest.RemoveStageFromRequest(req.StageRequestID, req.StageID, userID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove stage from request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Stage removed from request successfully"})
}

func (h *RequestStageHandler) UpdateStageToRequestConnection(ctx *gin.Context) {
	var req UpdateStageToRequestConnection
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage request data"})
		return
	}

	userID := GetFixedUserID()
	if err := h.repo.StageRequest.UpdateRequestToStage(req.StageRequestID, req.StageID, userID, &req.InputField1, &req.InputField2); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request stage"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request stage updated successfully"})
}
