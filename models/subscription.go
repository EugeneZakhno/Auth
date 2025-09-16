package models

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// Subscription represents a user's subscription to a service
type Subscription struct {
	ID         int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateSubscriptionRequest represents the request body for creating a subscription
type CreateSubscriptionRequest struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       int       `json:"price" binding:"required,min=1"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required"`
	EndDate     *string   `json:"end_date,omitempty"`
}

// UpdateSubscriptionRequest represents the request body for updating a subscription
type UpdateSubscriptionRequest struct {
	ServiceName string    `json:"service_name,omitempty"`
	Price       *int      `json:"price,omitempty" binding:"omitempty,min=1"`
	UserID      uuid.UUID `json:"user_id,omitempty"`
	StartDate   string    `json:"start_date,omitempty"`
	EndDate     *string   `json:"end_date,omitempty"`
}

// CalculateCostRequest represents the request for calculating total subscription cost
type CalculateCostRequest struct {
	UserID      *uuid.UUID `form:"user_id,omitempty"`
	ServiceName *string    `form:"service_name,omitempty"`
	StartPeriod string     `form:"start_period" binding:"required"`
	EndPeriod   string     `form:"end_period" binding:"required"`
}

// CalculateCostResponse represents the response for calculating total subscription cost
type CalculateCostResponse struct {
	TotalCost int `json:"total_cost"`
}

// Validate validates the subscription data
func (s *CreateSubscriptionRequest) Validate() error {
	// Validate date format (MM-YYYY)
	datePattern := regexp.MustCompile(`^(0[1-9]|1[0-2])-(\d{4})$`)
	if !datePattern.MatchString(s.StartDate) {
		return errors.New("start_date must be in MM-YYYY format")
	}

	if s.EndDate != nil && !datePattern.MatchString(*s.EndDate) {
		return errors.New("end_date must be in MM-YYYY format")
	}

	return nil
}

// Validate validates the subscription update data
func (s *UpdateSubscriptionRequest) Validate() error {
	// Validate date format (MM-YYYY)
	datePattern := regexp.MustCompile(`^(0[1-9]|1[0-2])-(\d{4})$`)
	if s.StartDate != "" && !datePattern.MatchString(s.StartDate) {
		return errors.New("start_date must be in MM-YYYY format")
	}

	if s.EndDate != nil && !datePattern.MatchString(*s.EndDate) {
		return errors.New("end_date must be in MM-YYYY format")
	}

	return nil
}

// Validate validates the calculate cost request
func (c *CalculateCostRequest) Validate() error {
	// Validate date format (MM-YYYY)
	datePattern := regexp.MustCompile(`^(0[1-9]|1[0-2])-(\d{4})$`)
	if !datePattern.MatchString(c.StartPeriod) {
		return errors.New("start_period must be in MM-YYYY format")
	}

	if !datePattern.MatchString(c.EndPeriod) {
		return errors.New("end_period must be in MM-YYYY format")
	}

	return nil
}