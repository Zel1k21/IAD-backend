package handler

import (
	"github.com/gin-gonic/gin"

	"iad-backend/internal/app/repository"
)

func RegisterHandlers(router *gin.Engine, repo *repository.Repository) {
	apiRouter := router.Group("/api")

	userHandler := NewUserHandler(repo)

	// Public routes (no auth required)
	publicRouter := apiRouter.Group("")
	{
		publicRouter.POST("/users/register", userHandler.Register)
		publicRouter.POST("/users/login", userHandler.Login)
		publicRouter.POST("/users/refresh", userHandler.RefreshToken)

		// Public stage routes
		stageHandler := NewStageHandler(repo)
		publicRouter.GET("/stages", stageHandler.GetStages)
		publicRouter.GET("/stages/:id", stageHandler.GetStageByID)
	}

	protectedRouter := apiRouter.Group("")
	protectedRouter.Use(userHandler.AuthMiddleware())
	{
		protectedRouter.GET("/users/profile", userHandler.GetProfile)
		protectedRouter.PUT("/users/profile", userHandler.UpdateProfile)
		protectedRouter.POST("/users/logout", userHandler.Logout)

		// Stage routes with role-based permissions
		stageHandler := NewStageHandler(repo)
		stageRouter := protectedRouter.Group("/stages")
		{
			stageRouter.POST("", userHandler.ScopeMiddleware("create:stages"), stageHandler.CreateStage)
			stageRouter.PUT("/:id", userHandler.ScopeMiddleware("update:stages"), stageHandler.UpdateStage)
			stageRouter.DELETE("/:id", userHandler.ScopeMiddleware("delete:stages"), stageHandler.DeleteStage)
			stageRouter.POST("/:id/image", userHandler.ScopeMiddleware("update:stages"), stageHandler.AddStageImage)
			stageRouter.POST("/:id/add-to-request", userHandler.ScopeMiddleware("create:requests"), stageHandler.AddStageToDraftRequest)
		}

		requestHandler := NewStageRequestHandler(repo)
		requestRouter := protectedRouter.Group("/stage-requests")
		{
			requestRouter.GET("/stageRequestInfo", requestHandler.GetStageRequestInfo)
			requestRouter.GET("", requestHandler.GetStageRequests)
			requestRouter.GET("/:id", requestHandler.GetStageRequestByID)
			requestRouter.PUT("/:id", userHandler.ScopeMiddleware("update:requests"), requestHandler.UpdateStageRequest)
			requestRouter.PUT("/:id/form", userHandler.ScopeMiddleware("update:requests"), requestHandler.FormStageRequest)
			requestRouter.PUT("/:id/resolve", userHandler.ScopeMiddleware("resolve:requests"), requestHandler.ResolveStageRequest)
			requestRouter.PUT("/:id/reject", userHandler.ScopeMiddleware("reject:requests"), requestHandler.RejectStageRequest)
			requestRouter.DELETE("/:id", userHandler.ScopeMiddleware("update:requests"), requestHandler.DeleteStageRequest)
		}

		requestStageHandler := NewStageRequestToStageHandler(repo)
		requestStageRouter := protectedRouter.Group("/stage-request-stages")
		{
			requestStageRouter.DELETE("/:requestId/stages/:stageId", userHandler.ScopeMiddleware("update:requests"), requestStageHandler.RemoveStageToRequestConnection)
			requestStageRouter.PUT("/:requestId/stages/:stageId", userHandler.ScopeMiddleware("update:requests"), requestStageHandler.UpdateStageToRequestConnection)
		}
	}
}
