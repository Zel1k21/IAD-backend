package handler

import (
	"iad-backend/internal/app/repository"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type StageRequestHandler struct {
	repo *repository.Repository
}

type StagesRequestsFilterResponse struct {
	ID          uint64    `json:"id"`
	Status      uint8     `json:"Status"`
	UserID      uint64    `json:"UserID"`
	ModeratorID uint64    `json:"ModeratorID"`
	CreatedAt   time.Time `json:"CreatedAt"`
	FormedAt    time.Time `json:"FormedAt"`
	ClosedAt    time.Time `json:"ClosedAt"`
	ProductName string    `json:"ProductName"`
}

type StageRequestResponse struct {
	ID                   uint64                        `json:"id"`
	ProductName          string                        `json:"ProductName"`
	StageRequestToStage  []StageRequestToStageResponse `json:"StageRequestToStage"`
	FirstDimensionName   string                        `json:"FirstDimensionName"`
	FirstDimensionConst  float64                       `json:"FirstDimensionConst"`
	SecondDimensionName  string                        `json:"SecondDimensionName"`
	SecondDimensionConst float64                       `json:"SecondDimensionConst"`
	InputField1          uint64                        `json:"input_field_1"`
	InputField2          uint64                        `json:"input_field_2"`
}

type StageRequestDetailResponse struct {
	ID                   uint64                              `json:"id"`
	CreatedAt            time.Time                           `json:"created_at"`
	ProductName          string                              `json:"product_name"`
	StageRequestToStages []StageRequestToStageDetailResponse `json:"stage_request_to_stages"`
}

type StageRequestToStageDetailResponse struct {
	FirstDimensionName   string  `json:"first_dimension_name"`
	FirstDimensionConst  float64 `json:"first_dimension_const"`
	SecondDimensionName  string  `json:"second_dimension_name"`
	SecondDimensionConst float64 `json:"second_dimension_const"`
	InputField1          uint64  `json:"input_field_1"`
	InputField2          uint64  `json:"input_field_2"`
	StageTitle           string  `json:"stage_title"`
}

func NewStageRequestHandler(repository *repository.Repository) *StageRequestHandler {
	return &StageRequestHandler{repo: repository}
}

// @Summary      Get draft request info
// @Description  Get information about current user's draft request
// @Tags         stage-requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  StageRequestInfoResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /stage-requests/stageRequestInfo [get]

type StageRequestInfoResponse struct {
	RequestID uint64 `json:"request_id"`
	ItemCount int    `json:"item_count"`
}

type UpdateStageRequestResponse struct {
	ProductName *string `json:"product_name"`
}

// @Summary      Get stage requests
// @Description  Get a list of stage requests with optional filtering
// @Tags         stage-requests
// @Accept       json
// @Produce      json
// @Param        status query int false "Filter by status"
// @Param        date_from query string false "Filter by date from (YYYY-MM-DD)"
// @Param        date_to query string false "Filter by date to (YYYY-MM-DD)"
// @Success      200  {array}   StagesRequestsFilterResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /stage-requests [get]
func (h *StageRequestHandler) GetStageRequestInfo(ctx *gin.Context) {
	userUUID, _, ok := GetUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusOK, StageRequestInfoResponse{
			RequestID: 0,
			ItemCount: -1,
		})
		return
	}

	user, err := h.repo.User.GetUserByUUID(userUUID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	requestID, itemCount, err := h.repo.StageRequest.GetDraftRequestInfo(user.ID)
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
	userUUID, scopes, ok := GetUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user ID from UUID
	user, err := h.repo.User.GetUserByUUID(userUUID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// Check if user is moderator (has any moderator scopes)
	isModerator := false
	for _, scope := range scopes {
		if scope == "resolve:requests" || scope == "reject:requests" || scope == "manage:users" {
			isModerator = true
			break
		}
	}

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

	requests, err := h.repo.StageRequest.GetStageRequests(user.ID, isModerator, statusFilter, dateFrom, dateTo)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stage requests"})
		return
	}

	var response []StagesRequestsFilterResponse
	for _, stageRequest := range requests {
		response = append(response, StagesRequestsFilterResponse{
			ID:          stageRequest.ID,
			Status:      stageRequest.Status,
			UserID:      stageRequest.UserID,
			ModeratorID: stageRequest.ModeratorID,
			CreatedAt:   stageRequest.CreatedAt,
			FormedAt:    stageRequest.FormedAt,
			ClosedAt:    stageRequest.ClosedAt,
			ProductName: stageRequest.ProductName,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary      Get stage request by ID
// @Description  Get detailed information about a specific stage request
// @Tags         stage-requests
// @Accept       json
// @Produce      json
// @Param        id path int true "Stage Request ID"
// @Security     BearerAuth
// @Success      200  {object}  StageRequestDetailResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /stage-requests/{id} [get]
func (h *StageRequestHandler) GetStageRequestByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	userUUID, scopes, ok := GetUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user ID from UUID
	user, err := h.repo.User.GetUserByUUID(userUUID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// Check if user is moderator (has any moderator scopes)
	isModerator := false
	for _, scope := range scopes {
		if scope == "resolve:requests" || scope == "reject:requests" || scope == "manage:users" {
			isModerator = true
			break
		}
	}

	request, err := h.repo.StageRequest.GetStageRequestByID(id, user.ID, isModerator)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Stage request not found"})
		return
	}

	response := StageRequestDetailResponse{
		ID:          request.ID,
		CreatedAt:   request.CreatedAt,
		ProductName: request.ProductName,
	}

	for _, stageToRequest := range request.StageRequestToStage {
		stageDetail := StageRequestToStageDetailResponse{
			StageTitle:           stageToRequest.Stage.Title,
			FirstDimensionName:   stageToRequest.Stage.FirstDimensionName,
			FirstDimensionConst:  stageToRequest.Stage.FirstDimensionConst,
			SecondDimensionName:  stageToRequest.Stage.SecondDimensionName,
			SecondDimensionConst: stageToRequest.Stage.SecondDimensionConst,
			InputField1:          stageToRequest.InputField1,
			InputField2:          stageToRequest.InputField2,
		}
		response.StageRequestToStages = append(response.StageRequestToStages, stageDetail)
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary      Update stage request
// @Description  Update an existing stage request
// @Tags         stage-requests
// @Accept       json
// @Produce      json
// @Param        id path int true "Stage Request ID"
// @Param        request body UpdateStageRequestResponse true "Stage request update data"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /stage-requests/{id} [put]
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

// @Summary      Form stage request
// @Description  Form a draft stage request into a submitted request
// @Tags         stage-requests
// @Accept       json
// @Produce      json
// @Param        id path int true "Stage Request ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /stage-requests/{id}/form [put]
func (h *StageRequestHandler) FormStageRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage request ID"})
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

	if err := h.repo.StageRequest.FormRequest(id, user.ID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Stage request formed successfully"})
}

// @Summary      Resolve stage request
// @Description  Resolve a stage request (moderator action)
// @Tags         stage-requests
// @Accept       json
// @Produce      json
// @Param        id path int true "Stage Request ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /stage-requests/{id}/resolve [put]
func (h *StageRequestHandler) ResolveStageRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage request ID"})
		return
	}

	userUUID, scopes, ok := GetUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	hasScope := slices.Contains(scopes, "resolve:requests")
	if !hasScope {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions. Moderator role required"})
		return
	}

	user, err := h.repo.User.GetUserByUUID(userUUID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	deliveryDate := time.Now().AddDate(0, 1, 0)

	emissionCalculationResult, err := h.repo.StageRequest.ResolveOrRejectRequest(id, user.ID, 4)
	if err != nil {
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

// @Summary      Reject stage request
// @Description  Reject a stage request (moderator action)
// @Tags         stage-requests
// @Accept       json
// @Produce      json
// @Param        id path int true "Stage Request ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /stage-requests/{id}/reject [put]
func (h *StageRequestHandler) RejectStageRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage request ID"})
		return
	}

	userUUID, scopes, ok := GetUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	hasScope := false
	for _, scope := range scopes {
		if scope == "reject:requests" {
			hasScope = true
			break
		}
	}
	if !hasScope {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions. Moderator role required"})
		return
	}

	user, err := h.repo.User.GetUserByUUID(userUUID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	_, err = h.repo.StageRequest.ResolveOrRejectRequest(id, user.ID, 5)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Stage request rejected successfully"})
}

// @Summary      Delete stage request
// @Description  Delete a stage request
// @Tags         stage-requests
// @Accept       json
// @Produce      json
// @Param        id path int true "Stage Request ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /stage-requests/{id} [delete]
func (h *StageRequestHandler) DeleteStageRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage request ID"})
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

	if err := h.repo.StageRequest.DeleteRequest(id, user.ID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete stage request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Stage request deleted successfully"})
}
