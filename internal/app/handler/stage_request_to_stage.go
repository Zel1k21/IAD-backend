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

// @Summary      Create new stage request to stage handler
// @Description  Initialize handler for managing stage-request relationships
// @Tags         stage-request-stages
func NewStageRequestToStageHandler(repo *repository.Repository) *RequestStageHandler {
	return &RequestStageHandler{
		repo: repo,
	}
}

type StageRequestToStageResponse struct {
	Stage       []StageResponse `json:"stage"`
	InputField1 uint64          `json:"input_field_1"`
	InputField2 uint64          `json:"input_field_2"`
}

type UpdateStageToRequestConnection struct {
	InputField1 *uint64 `json:"input_field_1"`
	InputField2 *uint64 `json:"input_field_2"`
}

// @Summary      Remove stage from request
// @Description  Remove a stage from a stage request
// @Tags         stage-request-stages
// @Accept       json
// @Produce      json
// @Param        requestId path int true "Stage Request ID"
// @Param        stageId path int true "Stage ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /stage-request-stages/{requestId}/stages/{stageId} [delete]
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

	userUUID, _, ok := GetUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.repo.User.GetUserByUUID(userUUID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	if err := h.repo.StageRequest.RemoveStageFromRequest(requestID, stageID, user.ID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove stage from request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Stage removed from request successfully"})
}

// @Summary      Update stage in request
// @Description  Update stage connection details in a stage request
// @Tags         stage-request-stages
// @Accept       json
// @Produce      json
// @Param        requestId path int true "Stage Request ID"
// @Param        stageId path int true "Stage ID"
// @Param        request body UpdateStageToRequestConnection true "Stage connection update data"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /stage-request-stages/{requestId}/stages/{stageId} [put]
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

	userUUID, _, ok := GetUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.repo.User.GetUserByUUID(userUUID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	if err := h.repo.StageRequest.UpdateRequestToStage(requestID, stageID, user.ID, req.InputField1, req.InputField2); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request stage"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request stage updated successfully"})
}
