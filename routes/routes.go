package routes

import (
	"product-tracker/handlers"
	"product-tracker/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	product := r.Group("/product")
	product.Use(middlewares.AuthMiddleware())
	{
		product.POST("/insert", handlers.ImportProduct)
		product.POST("/list", handlers.GetProducts)
		product.GET("/list/:name", handlers.GetProductsByName)
	}
	
	r.GET("/health", handlers.HealthCheck)

	return r
}