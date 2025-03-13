package routes

import (
	"fmt"
	"net/http"
	"product-tracker/config"
	"product-tracker/handlers"
	"product-tracker/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// SetupRouter initializes and configures the router with all necessary middleware and routes
func SetupRouter() *gin.Engine {
	// Set gin mode based on environment
	if config.GetConfig().Server.Port == "8080" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router with custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validators here
		_ = v
	}

	r := gin.New()

	// Add middleware
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s | %s | %s | %s | %s | %d | %s | %s\n",
			param.TimeStamp.Format(time.RFC3339),
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.UserAgent(),
			param.StatusCode,
			param.Latency,
			param.ErrorMessage,
		)
	}))
	r.Use(requestid.New())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check route
	r.GET("/health", handlers.HealthCheck)

	// API version group
	v1 := r.Group("/api/v1")
	{
		// Product routes
		product := v1.Group("/product")
		product.Use(middlewares.AuthMiddleware())
		{
			product.POST("/insert", handlers.ImportProduct)
			product.POST("/list", handlers.GetProducts)
			product.GET("/list/:name", handlers.GetProductsByName)
		}
	}

	// Handle 404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
	})

	return r
}
