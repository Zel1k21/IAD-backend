package handler

import (
	"iad-backend/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type StagesHandler struct {
	StageRepositiry       *repository.StageRepository
	CalcRequestRepository *repository.CalcRequestRepository
}

func NewStagesHandler(stageRepository *repository.StageRepository, calcRequestRepository *repository.CalcRequestRepository) *StagesHandler {
	return &StagesHandler{
		StageRepositiry:       stageRepository,
		CalcRequestRepository: calcRequestRepository,
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

	calcRequestID := 1
	calcRequestEntryCount, err := h.CalcRequestRepository.GetCalcRequestEntryCountByID(calcRequestID)
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
		"calcRequestID":         calcRequestID,
		"calcRequestEntryCount": calcRequestEntryCount,
	})
}
