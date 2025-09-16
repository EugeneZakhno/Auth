package repository

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"

	"subscription-service/models"
)

// MockSubscriptionRepository is a mock implementation of the SubscriptionRepository for testing
type MockSubscriptionRepository struct {
	subscriptions map[string]models.Subscription
	mutex         sync.RWMutex
}

// NewMockSubscriptionRepository creates a new instance of MockSubscriptionRepository
func NewMockSubscriptionRepository() *MockSubscriptionRepository {
	return &MockSubscriptionRepository{
		subscriptions: make(map[string]models.Subscription),
	}
}

// Create adds a new subscription to the mock repository
func (r *MockSubscriptionRepository) Create(subscription models.Subscription) (models.Subscription, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Generate a new ID if not provided
	if subscription.ID == 0 {
		subscription.ID = len(r.subscriptions) + 1
	}

	// Set created and updated timestamps
	now := time.Now()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now

	// Store the subscription
	r.subscriptions[strconv.Itoa(subscription.ID)] = subscription

	return subscription, nil
}

// GetByID retrieves a subscription by its ID
func (r *MockSubscriptionRepository) GetByID(id string) (models.Subscription, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	subscription, exists := r.subscriptions[id]
	if !exists {
		return models.Subscription{}, errors.New("subscription not found")
	}

	return subscription, nil
}

// List returns all subscriptions with optional filtering
func (r *MockSubscriptionRepository) List(userID, serviceName string) ([]models.Subscription, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []models.Subscription

	for _, sub := range r.subscriptions {
		// Apply filters if provided
		if userID != "" {
			uid, err := uuid.Parse(userID)
			if err != nil {
				continue
			}
			if sub.UserID != uid {
				continue
			}
		}
		if serviceName != "" && sub.ServiceName != serviceName {
			continue
		}

		result = append(result, sub)
	}

	return result, nil
}

// Update modifies an existing subscription
func (r *MockSubscriptionRepository) Update(id string, subscription models.Subscription) (models.Subscription, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.subscriptions[id]
	if !exists {
		return models.Subscription{}, errors.New("subscription not found")
	}

	// Update the subscription
	idInt, _ := strconv.Atoi(id)
	subscription.ID = idInt
	subscription.UpdatedAt = time.Now()
	r.subscriptions[id] = subscription

	return subscription, nil
}

// Delete removes a subscription by its ID
func (r *MockSubscriptionRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.subscriptions[id]
	if !exists {
		return errors.New("subscription not found")
	}

	delete(r.subscriptions, id)
	return nil
}

// CalculateTotalCost calculates the total cost of subscriptions for a given period and filters
func (r *MockSubscriptionRepository) CalculateTotalCost(userID, serviceName, startPeriod, endPeriod string) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Parse start and end periods
	startDate, err := time.Parse("01-2006", startPeriod)
	if err != nil {
		return 0, errors.New("invalid start period format")
	}

	endDate, err := time.Parse("01-2006", endPeriod)
	if err != nil {
		return 0, errors.New("invalid end period format")
	}

	if startDate.After(endDate) {
		return 0, errors.New("start period cannot be after end period")
	}

	// Get subscriptions that match the filters
	subscriptions, err := r.List(userID, serviceName)
	if err != nil {
		return 0, err
	}

	// Calculate total cost
	totalCost := 0
	for _, sub := range subscriptions {
		// Parse subscription start date
		subStartDate, err := time.Parse("01-2006", sub.StartDate)
		if err != nil {
			continue
		}

		// Skip if subscription starts after the end period
		if subStartDate.After(endDate) {
			continue
		}

		// Adjust start date if subscription starts after the requested start period
		effectiveStartDate := startDate
		if subStartDate.After(startDate) {
			effectiveStartDate = subStartDate
		}

		// Calculate effective end date
		effectiveEndDate := endDate
		if sub.EndDate != nil {
			subEndDate, err := time.Parse("01-2006", *sub.EndDate)
			if err == nil && subEndDate.Before(endDate) {
				effectiveEndDate = subEndDate
			}
		}

		// Calculate months for this subscription
		subMonths := (effectiveEndDate.Year()-effectiveStartDate.Year())*12 + int(effectiveEndDate.Month()-effectiveStartDate.Month()) + 1
		if subMonths > 0 {
			totalCost += sub.Price * subMonths
		}
	}

	return totalCost, nil
}
