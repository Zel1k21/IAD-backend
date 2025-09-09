package api

import (
	"iad-backend/internal/app/handler"
	"iad-backend/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	stageRepo, err := repository.NewStageRepository()
	if err != nil {
		logrus.Error("inicialize stage repository error: ", err)
	}

	stageHandler := handler.NewStageHandler(stageRepo)
	stagesHandler := handler.NewStagesHandler(stageRepo)

	r := gin.Default()

	r.LoadHTMLGlob("./templates/**/*")
	r.Static("/static", "./resources")

	r.GET("/stage/:id", stageHandler.GetStageByID)
	r.GET("/stages", stagesHandler.GetStages)

	r.Run()
	log.Println("server down")
}
