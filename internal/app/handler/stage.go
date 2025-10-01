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

func (h *StageHandler) GetStages(ctx *gin.Context) {
	searchQuery := ctx.Query("title")

	stage, err := h.repo.Stage.GetStages(searchQuery)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stages"})
		return
	}

	ctx.JSON(http.StatusOK, stage)
}

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

	ctx.JSON(http.StatusOK, stage)
}

func (h *StageHandler) CreateStage(ctx *gin.Context) {
	var req CreateStageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	ctx.JSON(http.StatusCreated, stage)
}

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

func (h *StageHandler) AddStageToDraftRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage ID"})
		return
	}

	userID := GetFixedUserID()
	if err := h.repo.Stage.AddStageToDraftRequest(id, userID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add stage to draft request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Stage added to draft request successfully"})
}
