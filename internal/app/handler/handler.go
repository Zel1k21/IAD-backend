package handler

import (
	"github.com/gin-gonic/gin"

	"iad-backend/internal/app/repository"
)

func RegisterHandlers(router *gin.Engine, repo *repository.Repository) {
	apiRouter := router.Group("/api")

	stageHandler := NewStageHandler(repo)
	stageRouter := apiRouter.Group("/stages")
	{
		stageRouter.GET("", stageHandler.GetStages)
		stageRouter.GET("/:id", stageHandler.GetStageByID)
		stageRouter.POST("", stageHandler.CreateStage)
		stageRouter.PUT("/:id", stageHandler.UpdateStage)
		stageRouter.DELETE("/:id", stageHandler.DeleteStage)
		stageRouter.POST("/:id/image", stageHandler.AddStageImage)
		stageRouter.POST("/:id/draft", stageHandler.AddStageToDraftRequest)
	}

	requestHandler := NewStageRequestHandler(repo)
	requestRouter := apiRouter.Group("/stage-requests")
	{
		requestRouter.GET("/cart", requestHandler.GetStageRequestInfo)
		requestRouter.GET("", requestHandler.GetStageRequests)
		requestRouter.GET("/:id", requestHandler.GetStageRequestByID)
		requestRouter.PUT("/:id", requestHandler.UpdateStageRequest)
		requestRouter.PUT("/:id/form", requestHandler.FormStageRequest)
		requestRouter.PUT("/:id/resolve", requestHandler.ResolveStageRequest)
		requestRouter.PUT("/:id/reject", requestHandler.RejectStageRequest)
		requestRouter.DELETE("/:id", requestHandler.DeleteStageRequest)
	}

	requestStageHandler := NewStageRequestToStageHandler(repo)
	requestStageRouter := apiRouter.Group("/light-request-stages")
	{
		requestStageRouter.DELETE("", requestStageHandler.RemoveStageToRequestConnection)
		requestStageRouter.PUT("", requestStageHandler.UpdateStageToRequestConnection)
	}

	userHandler := NewUserHandler(repo)
	userRouter := apiRouter.Group("/users")
	{
		userRouter.POST("/register", userHandler.Register)
		userRouter.GET("/profile", userHandler.GetProfile)
		userRouter.PUT("/profile", userHandler.UpdateProfile)
		userRouter.POST("/login", userHandler.Login)
		userRouter.POST("/logout", userHandler.Logout)
	}
}
