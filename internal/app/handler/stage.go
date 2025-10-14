package handler

import (
	"iad-backend/internal/app/ds"
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

type CreateStageRequest struct {
	Title                string  `json:"title" binding:"required"`
	Description          string  `json:"description"`
	FirstDimensionName   string  `json:"first_dimension_name" binding:"required"`
	SecondDimensionName  string  `json:"second_dimension_name" binding:"required"`
	FirstDimensionConst  float64 `json:"first_dimension_const" binding:"required"`
	SecondDimensionConst float64 `json:"second_dimension_const" binding:"required"`
}

type UpdateStageRequest struct {
	Title                *string  `json:"title"`
	Description          *string  `json:"description"`
	FirstDimensionName   *string  `json:"first_dimension_name"`
	SecondDimensionName  *string  `json:"second_dimension_name"`
	FirstDimensionConst  *float64 `json:"first_dimension_const"`
	SecondDimensionConst *float64 `json:"second_dimension_const"`
}

type StagesFilterResponse struct {
	ID       uint64 `json:"id"`
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
}

type StageResponse struct {
	ID                   uint64  `json:"id"`
	Title                string  `json:"title"`
	ImageURL             string  `json:"image_url"`
	Description          string  `json:"description"`
	FirstDimensionName   string  `json:"first_dimension_name"`
	FirstDimensionConst  float64 `json:"first_dimension_const"`
	SecondDimensionName  string  `json:"second_dimension_name"`
	SecondDimensionConst float64 `json:"second_dimension_const"`
}

// @Summary      Get all stages
// @Description  Get a list of all stages with optional title search
// @Tags         stages
// @Accept       json
// @Produce      json
// @Param        title query string false "Search stages by title"

// @Success      200  {array}   StagesFilterResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /stages [get]
func (h *StageHandler) GetStages(ctx *gin.Context) {
	searchQuery := ctx.Query("title")

	stages, err := h.repo.Stage.GetStages(searchQuery)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stages"})
		return
	}

	var response []StagesFilterResponse
	for _, stage := range stages {
		response = append(response, StagesFilterResponse{
			ID:       stage.ID,
			Title:    stage.Title,
			ImageURL: stage.ImageURL,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary      Get stage by ID
// @Description  Get detailed information about a specific stage
// @Tags         stages
// @Accept       json
// @Produce      json
// @Param        id path int true "Stage ID"

// @Success      200  {object}  StageResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /stages/{id} [get]
func (h *StageHandler) GetStageByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage ID"})
		return
	}

	stage, err := h.repo.Stage.GetStageByID(id)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Stage not found"})
		return
	}

	responce := StageResponse{
		ID:          stage.ID,
		Title:       stage.Title,
		ImageURL:    stage.ImageURL,
		Description: stage.Description,
	}
	ctx.JSON(http.StatusOK, responce)
}

// @Summary      Create a new stage
// @Description  Create a new stage with the provided data
// @Tags         stages
// @Accept       json
// @Produce      json
// @Param        request body CreateStageRequest true "Stage creation data"
// @Security     BearerAuth
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /stages [post]
func (h *StageHandler) CreateStage(ctx *gin.Context) {
	var req CreateStageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("CreateStage validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	stage := &ds.Stage{
		Title:                req.Title,
		Description:          req.Description,
		FirstDimensionName:   req.FirstDimensionName,
		SecondDimensionName:  req.SecondDimensionName,
		FirstDimensionConst:  req.FirstDimensionConst,
		SecondDimensionConst: req.SecondDimensionConst,
	}

	if err := h.repo.Stage.CreateStage(stage); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stage"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "Stage created successfully",
		"stage_id": stage.ID,
	})
}

// @Summary      Update stage
// @Description  Update an existing stage with new data
// @Tags         stages
// @Accept       json
// @Produce      json
// @Param        id path int true "Stage ID"
// @Param        request body UpdateStageRequest true "Stage update data"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /stages/{id} [put]
func (h *StageHandler) UpdateStage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage ID"})
		return
	}

	var req UpdateStageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	stageData := &ds.Stage{}
	if req.Title != nil {
		stageData.Title = *req.Title
	}
	if req.Description != nil {
		stageData.Description = *req.Description
	}
	if req.FirstDimensionName != nil {
		stageData.FirstDimensionName = *req.FirstDimensionName
	}
	if req.SecondDimensionName != nil {
		stageData.SecondDimensionName = *req.SecondDimensionName
	}
	if req.FirstDimensionConst != nil {
		stageData.FirstDimensionConst = *req.FirstDimensionConst
	}
	if req.SecondDimensionConst != nil {
		stageData.SecondDimensionConst = *req.SecondDimensionConst
	}

	if err := h.repo.Stage.UpdateStage(id, stageData); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stage"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Stage updated successfully"})
}

// @Summary      Delete stage
// @Description  Delete a stage by ID
// @Tags         stages
// @Accept       json
// @Produce      json
// @Param        id path int true "Stage ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /stages/{id} [delete]
func (h *StageHandler) DeleteStage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage ID"})
		return
	}

	if err := h.repo.Stage.DeleteStage(id); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete stage"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Stage deleted successfully"})
}

// @Summary      Add image to stage
// @Description  Upload and attach an image to a stage
// @Tags         stages
// @Accept       multipart/form-data
// @Produce      json
// @Param        id path int true "Stage ID"
// @Param        image formData file true "Stage image file"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /stages/{id}/image [post]
func (h *StageHandler) AddStageImage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage ID"})
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	if err := h.repo.Stage.AddStageImage(id, file); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add image"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Image added successfully"})
}

// @Summary      Add stage to draft request
// @Description  Add a stage to the current user's draft request
// @Tags         stages
// @Accept       json
// @Produce      json
// @Param        id path int true "Stage ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /stages/{id}/add-to-request [post]
func (h *StageHandler) AddStageToDraftRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		logrus.Errorf("AddStageToDraftRequest validation failed: %v", err)
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

	if err := h.repo.Stage.AddStageToDraftRequest(id, user.ID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add stage to draft request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Stage added to draft request successfully"})
}
