package main

import (
	"fmt"
	"log"
	"net/http"
	"subscription-service/config"
	"subscription-service/db"
	"subscription-service/handlers"
	"subscription-service/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "subscription-service/docs"
)

// @title Subscription Service API
// @version 1.0
// @description API for managing subscription services
// @host localhost:8081
// @BasePath /api/v1
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logger.NewLogger()

	// Connect to database
	postgres, err := db.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer postgres.Close()

	// Run database migrations
	if err := postgres.RunMigrations(); err != nil {
		logger.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize handlers with DB
	subscriptionHandler := handlers.NewSubscriptionHandler(postgres, logger)

	// Initialize router
	router := gin.Default()

	// Setup API routes
	api := router.Group("/api/v1")
	{
		// Health check endpoint
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "healthy",
				"mode":   "database",
			})
		})

		// Real subscriptions endpoints
		api.POST("/subscriptions", subscriptionHandler.Create)
		api.GET("/subscriptions/:id", subscriptionHandler.Get)
		api.GET("/subscriptions", subscriptionHandler.List)
		api.PUT("/subscriptions/:id", subscriptionHandler.Update)
		api.DELETE("/subscriptions/:id", subscriptionHandler.Delete)
		api.GET("/subscriptions/calculate", subscriptionHandler.CalculateTotalCost)
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, 8081)
	logger.Infof("Starting server on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
