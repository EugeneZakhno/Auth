package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"subscription-service/config"
	"subscription-service/handlers"
	"subscription-service/logger"
	"subscription-service/models"
	"subscription-service/repository"
)

func setupTestRouter() *gin.Engine {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a mock repository for testing
	repo := repository.NewMockSubscriptionRepository()

	// Create a new logger
	log := logger.NewLogger()

	// Create a new handler with the mock repository
	handler := handlers.NewSubscriptionHandler(repo, log)

	// Setup router
	r := gin.Default()

	// Register routes
	v1 := r.Group("/api/v1")
	{
		v1.POST("/subscriptions", handler.CreateSubscription)
		v1.GET("/subscriptions/:id", handler.GetSubscription)
		v1.GET("/subscriptions", handler.ListSubscriptions)
		v1.PUT("/subscriptions/:id", handler.UpdateSubscription)
		v1.DELETE("/subscriptions/:id", handler.DeleteSubscription)
		v1.GET("/subscriptions/cost", handler.CalculateTotalCost)
	}

	return r
}

func TestCreateSubscription(t *testing.T) {
	r := setupTestRouter()

	// Create a valid subscription
	userID := uuid.New()
	subscription := models.CreateSubscriptionRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      userID,
		StartDate:   "07-2025",
	}

	jsonValue, _ := json.Marshal(subscription)
	req, _ := http.NewRequest("POST", "/api/v1/subscriptions", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse the response
	var response models.SubscriptionResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check the response
	assert.Equal(t, subscription.ServiceName, response.ServiceName)
	assert.Equal(t, subscription.Price, response.Price)
	assert.Equal(t, subscription.UserID.String(), response.UserID)
	assert.Equal(t, subscription.StartDate, response.StartDate)
}

func TestGetSubscription(t *testing.T) {
	r := setupTestRouter()

	// First create a subscription
	userID := uuid.New()
	subscription := models.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       700,
		UserID:      userID,
		StartDate:   "01-2024",
	}

	jsonValue, _ := json.Marshal(subscription)
	req, _ := http.NewRequest("POST", "/api/v1/subscriptions", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var createResponse models.SubscriptionResponse
	json.Unmarshal(w.Body.Bytes(), &createResponse)

	// Now get the subscription
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/subscriptions/%s", createResponse.ID), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response
	var getResponse models.SubscriptionResponse
	err := json.Unmarshal(w.Body.Bytes(), &getResponse)
	assert.NoError(t, err)

	// Check the response
	assert.Equal(t, createResponse.ID, getResponse.ID)
	assert.Equal(t, subscription.ServiceName, getResponse.ServiceName)
	assert.Equal(t, subscription.Price, getResponse.Price)
	assert.Equal(t, subscription.UserID.String(), getResponse.UserID)
	assert.Equal(t, subscription.StartDate, getResponse.StartDate)
}

func TestListSubscriptions(t *testing.T) {
	r := setupTestRouter()

	// Create a few subscriptions
	userID := uuid.New()
	subscriptions := []models.CreateSubscriptionRequest{
		{
			ServiceName: "Spotify",
			Price:       199,
			UserID:      userID,
			StartDate:   "03-2024",
		},
		{
			ServiceName: "YouTube Premium",
			Price:       299,
			UserID:      userID,
			StartDate:   "04-2024",
		},
	}

	for _, sub := range subscriptions {
		jsonValue, _ := json.Marshal(sub)
		req, _ := http.NewRequest("POST", "/api/v1/subscriptions", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	// Now list all subscriptions
	req, _ := http.NewRequest("GET", "/api/v1/subscriptions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response
	var listResponse []models.SubscriptionResponse
	err := json.Unmarshal(w.Body.Bytes(), &listResponse)
	assert.NoError(t, err)

	// Check that we have at least the number of subscriptions we created
	assert.GreaterOrEqual(t, len(listResponse), len(subscriptions))
}

func TestCalculateTotalCost(t *testing.T) {
	r := setupTestRouter()

	// Create subscriptions for two different users
	userID1 := uuid.New()
	userID2 := uuid.New()

	subscriptions := []models.CreateSubscriptionRequest{
		{
			ServiceName: "Netflix",
			Price:       700,
			UserID:      userID1,
			StartDate:   "01-2024",
		},
		{
			ServiceName: "Spotify",
			Price:       199,
			UserID:      userID1,
			StartDate:   "02-2024",
		},
		{
			ServiceName: "YouTube Premium",
			Price:       299,
			UserID:      userID2,
			StartDate:   "03-2024",
		},
	}

	for _, sub := range subscriptions {
		jsonValue, _ := json.Marshal(sub)
		req, _ := http.NewRequest("POST", "/api/v1/subscriptions", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	// Calculate total cost for userID1
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/subscriptions/cost?user_id=%s&start_period=01-2024&end_period=12-2024", userID1.String()), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response
	var costResponse models.CalculateTotalCostResponse
	err := json.Unmarshal(w.Body.Bytes(), &costResponse)
	assert.NoError(t, err)

	// Check the total cost (700 + 199) * 12 = 10788
	assert.Equal(t, 10788, costResponse.TotalCost)
}