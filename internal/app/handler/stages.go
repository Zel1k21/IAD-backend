package handler

import (
	"iad-backend/internal/app/ds"
	"iad-backend/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type StagesHandler struct {
	repo *repository.Repository
}

func NewStagesHandler(repository *repository.Repository) *StagesHandler {
	return &StagesHandler{repo: repository}
}

func (h *StagesHandler) Register(router *gin.Engine) {
	router.GET("/stages", h.GetStages)
	router.POST("/stages", h.AddStageToRequest)
}

type CompTextInput struct {
	ShowLabel   bool
	Label       string
	Type        string
	Name        string
	Placeholder string
	Value       string
}

func (h *StagesHandler) GetStages(ctx *gin.Context) {
	var stages []ds.Stage
	var err error

	searchQuery := ctx.Query("title")
	if searchQuery == "" {
		stages, err = h.repo.Stage.GetStages()
		if err != nil {
			logrus.Error(err)
			ctx.Status(http.StatusNotFound)
			return
		}
	} else {
		stages, err = h.repo.Stage.GetStagesByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
			ctx.Status(http.StatusNotFound)
			return
		}
	}

	stageRequestID, stageRequestEntryCount, err := h.repo.StageRequest.GetStageRequestIDEntryCountByUserID(1)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "stages.html", gin.H{
		"title":  "Этапы жизненного цикла",
		"stages": stages,
		"search": CompTextInput{
			ShowLabel:   true,
			Label:       "Поиск",
			Type:        "text",
			Name:        "title",
			Placeholder: "Введите название этапа",
			Value:       searchQuery,
		},
		"stageRequestID":         stageRequestID,
		"stageRequestEntryCount": stageRequestEntryCount,
	})
}

func (h *StagesHandler) AddStageToRequest(ctx *gin.Context) {
	stageIDStr := ctx.PostForm("stage-id")
	stageID, err := strconv.ParseUint(stageIDStr, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	err = h.repo.StageRequest.AddStageToStageRequest(stageID, 1)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/stages")
}
