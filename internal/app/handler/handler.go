package handler

import (
	"github.com/gin-gonic/gin"

	"iad-backend/internal/app/repository"
)

func RegisterHandlers(router *gin.Engine, repo *repository.Repository) {
	router.LoadHTMLGlob("./templates/**/*")
	router.Static("/static", "./resources")

	StageRequestHandler := NewStageRequestHandler(repo)
	StageHandler := NewStageHandler(repo)
	StagesHandler := NewStagesHandler(repo)

	StageRequestHandler.Register(router)
	StageHandler.Register(router)
	StagesHandler.Register(router)
}
