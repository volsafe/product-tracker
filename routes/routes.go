package routes

import (
	"profile/handlers"
	"profile/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/profile")
	auth.Use(middlewares.AuthMiddleware())
	{
		auth.POST("/create", handlers.CreateProfile)
		auth.GET("/:userID", handlers.GetProfile)
		auth.PUT("/update", handlers.UpdateProfile)
		auth.DELETE("/delete/:userID", handlers.DeleteProfile)
	}
	
	r.GET("/health", handlers.HealthCheck)

	return r
}