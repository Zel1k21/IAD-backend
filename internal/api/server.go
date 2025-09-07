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

	repository, err := repository.NewRepository()
	if err != nil {
		logrus.Error("error creating repository")
	}
	handler := handler.NewHandler(repository)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/hello", handler.GetOrders)
	r.GET("/order/:id", handler.GetOrder)

	r.Run()
	log.Println("server down")
}
