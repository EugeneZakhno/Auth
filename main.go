package main

import (
	"log"
	"subscription-service/config"
	"subscription-service/db"
	"subscription-service/handlers"
	"subscription-service/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "subscription-service/docs" // This is for swagger
)

// @title           Subscription Service API
// @version         1.0
// @description     A REST service for aggregating data on users' online subscriptions
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	// Initialize logger
	logger := logger.NewLogger()
	logger.Info("Starting subscription service")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := db.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize router
	router := gin.Default()

	// Initialize API handlers
	subscriptionHandler := handlers.NewSubscriptionHandler(db, logger)

	// API routes
	api := router.Group("/api/v1")
	{
		subscriptions := api.Group("/subscriptions")
		{
			subscriptions.POST("", subscriptionHandler.Create)
			subscriptions.GET("", subscriptionHandler.List)
			subscriptions.GET("/:id", subscriptionHandler.Get)
			subscriptions.PUT("/:id", subscriptionHandler.Update)
			subscriptions.DELETE("/:id", subscriptionHandler.Delete)
			subscriptions.GET("/calculate", subscriptionHandler.CalculateTotalCost)
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	logger.Infof("Server starting on %s", serverAddr)
	if err := router.Run("0.0.0.0:" + cfg.Server.Port); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
