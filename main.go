package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize router
	router := gin.Default()

	// Setup API routes
	api := router.Group("/api/v1")
	{
		// Health check endpoint
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "healthy",
				"mode":   "local development",
			})
		})

		// Mock subscriptions endpoints
		api.GET("/subscriptions", func(c *gin.Context) {
			c.JSON(http.StatusOK, []gin.H{
				{
					"id":           1,
					"user_id":      "user123",
					"service_name": "Netflix",
					"cost":         9.99,
					"renewal_date": "2023-12-01",
				},
				{
					"id":           2,
					"user_id":      "user123",
					"service_name": "Spotify",
					"cost":         4.99,
					"renewal_date": "2023-12-15",
				},
			})
		})

		// Create subscription endpoint
		api.POST("/subscriptions", func(c *gin.Context) {
			var subscription map[string]interface{}
			if err := c.ShouldBindJSON(&subscription); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Just return the subscription with an ID
			subscription["id"] = 3
			c.JSON(http.StatusCreated, subscription)
		})
	}

	// Start the server on a different port to avoid conflicts
	log.Println("Starting server on localhost:8081")
	if err := router.Run("localhost:8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
