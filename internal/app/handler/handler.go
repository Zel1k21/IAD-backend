package handler

import (
	"github.com/gin-gonic/gin"

	"iad-backend/internal/app/repository"
)

func RegisterHandlers(router *gin.Engine, repo *repository.Repository) {
	apiRouter := router.Group("/api", SetCorsHeaders)

	userHandler := NewUserHandler(repo)

	// Public routes (no auth required)
	publicRouter := apiRouter.Group("")
	{
		publicRouter.POST("/users/register", userHandler.Register)
		publicRouter.OPTIONS("/users/register", SetCorsHeaders)
		publicRouter.POST("/users/login", userHandler.Login)
		publicRouter.OPTIONS("/users/login", SetCorsHeaders)
		publicRouter.POST("/users/refresh", userHandler.RefreshToken)
		publicRouter.OPTIONS("/users/refresh", SetCorsHeaders)

		// Public stage routes
		stageHandler := NewStageHandler(repo)
		publicRouter.GET("/stages", stageHandler.GetStages)
		publicRouter.GET("/stages/:id", stageHandler.GetStageByID)

		requestHandler := NewStageRequestHandler(repo)
		publicRouter.GET("/stage-requests/stageRequestInfo", requestHandler.GetStageRequestInfo)
		publicRouter.OPTIONS("/stage-requests/stageRequestInfo", SetCorsHeaders)
		publicRouter.PUT("/stage-requests/asyncUpdateCalculation", requestHandler.AsyncUpdateEmissionCalculation)
		publicRouter.OPTIONS("/stage-requests/asyncUpdateCalculation", SetCorsHeaders)

		publicRouter.OPTIONS("/users/profile", SetCorsHeaders)
		publicRouter.OPTIONS("/users/logout", SetCorsHeaders)
		publicRouter.OPTIONS("/stages", SetCorsHeaders)
		publicRouter.OPTIONS("/stages/:id", SetCorsHeaders)
		publicRouter.OPTIONS("/stages/:id/image", SetCorsHeaders)
		publicRouter.OPTIONS("/stages/:id/add-to-request", SetCorsHeaders)
		publicRouter.OPTIONS("/stage-requests", SetCorsHeaders)
		publicRouter.OPTIONS("/stage-requests/:id", SetCorsHeaders)
		publicRouter.OPTIONS("/stage-requests/:id/form", SetCorsHeaders)
		publicRouter.OPTIONS("/stage-requests/:id/resolve", SetCorsHeaders)
		publicRouter.OPTIONS("/stage-requests/:id/reject", SetCorsHeaders)
		publicRouter.OPTIONS("/stage-request-stages/:requestId/stages/:stageId", SetCorsHeaders)
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
