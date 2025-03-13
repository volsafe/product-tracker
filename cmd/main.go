package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"product-tracker/config"
	"product-tracker/routes"

	_ "product-tracker/docs" // Import swagger docs

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Product Tracker API
// @version         1.0
// @description     A RESTful API for tracking products and their energy consumption.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	config *config.Config
}

// NewServer creates a new HTTP server
func NewServer(router *gin.Engine, cfg *config.Config) *Server {
	return &Server{
		router: router,
		config: cfg,
	}
}

// NewRouter creates a new Gin router
func NewRouter(cfg *config.Config) *gin.Engine {
	// Set gin mode
	gin.SetMode(gin.DebugMode)

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup routes
	routes.SetupRoutes(router)

	return router
}

func main() {
	// Initialize logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("📝 Starting Product Tracker API...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Create router
	router := NewRouter(cfg)

	// Create server
	server := NewServer(router, cfg)

	// Create server address
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	serverURL := fmt.Sprintf("http://localhost%s", addr)

	// Create HTTP server
	srv := &http.Server{
		Addr:           addr,
		Handler:        server.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	log.Printf("🚀 Starting server on %s", serverURL)
	log.Printf("📚 Swagger documentation available at %s/swagger/index.html", serverURL)

	// Start server
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
