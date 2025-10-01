package handler

import (
	"iad-backend/internal/app/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type StageRequestHandler struct {
	repo *repository.Repository
}

func NewStageRequestHandler(repository *repository.Repository) *StageRequestHandler {
	return &StageRequestHandler{repo: repository}
}

type StageRequestInfoResponse struct {
	RequestID uint64 `json:"request_id"`
	ItemCount int    `json:"item_count"`
}

type UpdateStageRequestResponse struct {
	ProductName *string `json:"product_name"`
}

func (h *StageRequestHandler) GetStageRequestInfo(ctx *gin.Context) {
	userID := GetFixedUserID()
	requestID, itemCount, err := h.repo.StageRequest.GetDraftRequestInfo(userID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stage request info"})
		return
	}

	ctx.JSON(http.StatusOK, StageRequestInfoResponse{
		RequestID: requestID,
		ItemCount: itemCount,
	})
}

func (h *StageRequestHandler) GetStageRequests(ctx *gin.Context) {
	var statusFilter uint8
	if statusStr := ctx.Query("status"); statusStr != "" {
		if status, err := strconv.ParseUint(statusStr, 10, 8); err == nil {
			statusFilter = uint8(status)
		}
	}

	var dateFrom, dateTo *time.Time
	if dateFromStr := ctx.Query("date_from"); dateFromStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFrom = &parsed
		}
	}
	if dateToStr := ctx.Query("date_to"); dateToStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateTo = &parsed
		}
	}

	requests, err := h.repo.StageRequest.GetStageRequests(statusFilter, dateFrom, dateTo)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stage requests"})
		return
	}

	ctx.JSON(http.StatusOK, requests)
}

func (h *StageRequestHandler) GetStageRequestByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	userID := GetFixedUserID()
	request, err := h.repo.StageRequest.GetStageRequestByID(id, userID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Stage request not found"})
		return
	}

	ctx.JSON(http.StatusOK, request)
}

func (h *StageRequestHandler) UpdateStageRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage request ID"})
		return
	}

	var req UpdateStageRequestResponse
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage request data"})
		return
	}

	if err := h.repo.StageRequest.UpdateStageRequest(id, req.ProductName); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stage request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Stage request updated successfully"})
}

func (h *StageRequestHandler) FormStageRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage request ID"})
		return
	}

	userID := GetFixedUserID()
	if err := h.repo.StageRequest.FormRequest(id, userID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Stage request formed successfully"})
}

func (h *StageRequestHandler) ResolveStageRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage request ID"})
		return
	}

	emissionCalculationResult := h.repo.StageRequest.CalculateEmission(id)

	deliveryDate := time.Now().AddDate(0, 1, 0)

	moderatorID := uint64(2)
	if err := h.repo.StageRequest.ResolveOrRejectRequest(id, moderatorID, 4); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"message": "Stage request resolved successfully",
		"calculated_data": gin.H{
			"emission calculation result": emissionCalculationResult,
			"delivery_date":               deliveryDate.Format("2006-01-02"),
		},
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *StageRequestHandler) RejectStageRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage request ID"})
		return
	}

	moderatorID := uint64(2)
	if err := h.repo.StageRequest.ResolveOrRejectRequest(id, moderatorID, 5); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Stage request rejected successfully"})
}

func (h *StageRequestHandler) DeleteStageRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage request ID"})
		return
	}

	userID := GetFixedUserID()
	if err := h.repo.StageRequest.DeleteRequest(id, userID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete stage request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Stage request deleted successfully"})
}
