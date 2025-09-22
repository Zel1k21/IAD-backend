package handler

import (
	"iad-backend/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type StagesHandler struct {
	StageRepositiry        *repository.StageRepository
	StageRequestRepository *repository.StageRequestRepository
}

func NewStagesHandler(stageRepository *repository.StageRepository, stageRequestRepository *repository.StageRequestRepository) *StagesHandler {
	return &StagesHandler{
		StageRepositiry:        stageRepository,
		StageRequestRepository: stageRequestRepository,
	}
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
	var stages []repository.Stage
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		stages, err = h.StageRepositiry.GetStages()
		if err != nil {
			logrus.Error(err)
			ctx.Status(http.StatusNotFound)
			return
		}
	} else {
		stages, err = h.StageRepositiry.GetStagesByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
			ctx.Status(http.StatusNotFound)
			return
		}
	}

	stageRequestID := 1
	stageRequestEntryCount, err := h.StageRequestRepository.GetStageRequestEntryCountByID(stageRequestID)
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
			Type:        "text",
			Name:        "query",
			Placeholder: "Введите название этапа",
			Value:       searchQuery,
		},
		"stageRequestID":         stageRequestID,
		"stageRequestEntryCount": stageRequestEntryCount,
	})
}
