package handlers

import (
	"net/http"
	"strconv"
	"subscription-service/db"
	"subscription-service/logger"
	"subscription-service/models"
	"subscription-service/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SubscriptionHandler handles HTTP requests for subscriptions
type SubscriptionHandler struct {
	repo   *repository.SubscriptionRepository
	logger *logger.Logger
}

// NewSubscriptionHandler creates a new subscription handler
func NewSubscriptionHandler(db *db.PostgresDB, logger *logger.Logger) *SubscriptionHandler {
	repo := repository.NewSubscriptionRepository(db)
	return &SubscriptionHandler{repo: repo, logger: logger}
}

// Create godoc
// @Summary Create a new subscription
// @Description Create a new subscription record
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	var req models.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		h.logger.Errorf("Validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.repo.Create(&req)
	if err != nil {
		h.logger.Errorf("Failed to create subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	subscription, err := h.repo.GetByID(id)
	if err != nil {
		h.logger.Errorf("Failed to get created subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created subscription"})
		return
	}

	h.logger.Infof("Created subscription with ID: %d", id)
	c.JSON(http.StatusCreated, subscription)
}

// Get godoc
// @Summary Get a subscription
// @Description Get a subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorf("Invalid ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	subscription, err := h.repo.GetByID(id)
	if err != nil {
		h.logger.Errorf("Failed to get subscription: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	h.logger.Infof("Retrieved subscription with ID: %d", id)
	c.JSON(http.StatusOK, subscription)
}

// List godoc
// @Summary List subscriptions
// @Description List all subscriptions with optional filtering
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "Filter by user ID"
// @Param service_name query string false "Filter by service name"
// @Success 200 {array} models.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) List(c *gin.Context) {
	userIDStr := c.Query("user_id")
	serviceName := c.Query("service_name")

	var userID *uuid.UUID
	if userIDStr != "" {
		parsedID, err := uuid.Parse(userIDStr)
		if err != nil {
			h.logger.Errorf("Invalid user ID: %s", userIDStr)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
		userID = &parsedID
	}

	var serviceNamePtr *string
	if serviceName != "" {
		serviceNamePtr = &serviceName
	}

	subscriptions, err := h.repo.List(userID, serviceNamePtr)
	if err != nil {
		h.logger.Errorf("Failed to list subscriptions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list subscriptions"})
		return
	}

	h.logger.Infof("Listed %d subscriptions", len(subscriptions))
	c.JSON(http.StatusOK, subscriptions)
}

// Update godoc
// @Summary Update a subscription
// @Description Update a subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param subscription body models.UpdateSubscriptionRequest true "Subscription data"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorf("Invalid ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	var req models.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		h.logger.Errorf("Validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(id, &req); err != nil {
		h.logger.Errorf("Failed to update subscription: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	subscription, err := h.repo.GetByID(id)
	if err != nil {
		h.logger.Errorf("Failed to get updated subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated subscription"})
		return
	}

	h.logger.Infof("Updated subscription with ID: %d", id)
	c.JSON(http.StatusOK, subscription)
}

// Delete godoc
// @Summary Delete a subscription
// @Description Delete a subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorf("Invalid ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		h.logger.Errorf("Failed to delete subscription: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	h.logger.Infof("Deleted subscription with ID: %d", id)
	c.Status(http.StatusNoContent)
}

// CalculateTotalCost godoc
// @Summary Calculate total subscription cost
// @Description Calculate the total cost of subscriptions for a period
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "Filter by user ID"
// @Param service_name query string false "Filter by service name"
// @Param start_period query string true "Start period (MM-YYYY)"
// @Param end_period query string true "End period (MM-YYYY)"
// @Success 200 {object} models.CalculateCostResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/calculate [get]
func (h *SubscriptionHandler) CalculateTotalCost(c *gin.Context) {
	var req models.CalculateCostRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Errorf("Failed to bind query: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		h.logger.Errorf("Validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.Query("user_id")
	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			h.logger.Errorf("Invalid user ID: %s", userIDStr)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
		req.UserID = &userID
	}

	serviceName := c.Query("service_name")
	if serviceName != "" {
		req.ServiceName = &serviceName
	}

	totalCost, err := h.repo.CalculateTotalCost(&req)
	if err != nil {
		h.logger.Errorf("Failed to calculate total cost: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate total cost"})
		return
	}

	h.logger.Infof("Calculated total cost: %d", totalCost)
	c.JSON(http.StatusOK, models.CalculateCostResponse{TotalCost: totalCost})
}