package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"subscription-service/db"
	"subscription-service/models"
	"time"

	"github.com/google/uuid"
)

// SubscriptionRepository handles database operations for subscriptions
type SubscriptionRepository struct {
	db *db.PostgresDB
}

// NewSubscriptionRepository creates a new subscription repository
func NewSubscriptionRepository(db *db.PostgresDB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create creates a new subscription
func (r *SubscriptionRepository) Create(subscription *models.CreateSubscriptionRequest) (int, error) {
	var id int
	err := r.db.DB.QueryRow(
		`INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		subscription.ServiceName, subscription.Price, subscription.UserID, subscription.StartDate, subscription.EndDate,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create subscription: %w", err)
	}

	return id, nil
}

// GetByID gets a subscription by ID
func (r *SubscriptionRepository) GetByID(id int) (*models.Subscription, error) {
	var subscription models.Subscription
	var createdAt, updatedAt time.Time

	err := r.db.DB.QueryRow(
		`SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at 
		FROM subscriptions WHERE id = $1`,
		id,
	).Scan(
		&subscription.ID,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserID,
		&subscription.StartDate,
		&subscription.EndDate,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	subscription.CreatedAt = createdAt
	subscription.UpdatedAt = updatedAt

	return &subscription, nil
}

// List gets all subscriptions with optional filtering
func (r *SubscriptionRepository) List(userID *uuid.UUID, serviceName *string) ([]*models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at FROM subscriptions`

	whereConditions := []string{}
	args := []interface{}{}
	paramCounter := 1

	if userID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("user_id = $%d", paramCounter))
		args = append(args, *userID)
		paramCounter++
	}

	if serviceName != nil && *serviceName != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("service_name = $%d", paramCounter))
		args = append(args, *serviceName)
		paramCounter++
	}

	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	rows, err := r.db.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer rows.Close()

	subscriptions := []*models.Subscription{}
	for rows.Next() {
		var subscription models.Subscription
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&subscription.ID,
			&subscription.ServiceName,
			&subscription.Price,
			&subscription.UserID,
			&subscription.StartDate,
			&subscription.EndDate,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}

		subscription.CreatedAt = createdAt
		subscription.UpdatedAt = updatedAt
		subscriptions = append(subscriptions, &subscription)
	}

	return subscriptions, nil
}

// Update updates a subscription
func (r *SubscriptionRepository) Update(id int, subscription *models.UpdateSubscriptionRequest) error {
	setClauses := []string{}
	args := []interface{}{}
	paramCounter := 1

	if subscription.ServiceName != "" {
		setClauses = append(setClauses, fmt.Sprintf("service_name = $%d", paramCounter))
		args = append(args, subscription.ServiceName)
		paramCounter++
	}

	if subscription.Price != nil {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d", paramCounter))
		args = append(args, *subscription.Price)
		paramCounter++
	}

	if subscription.UserID != uuid.Nil {
		setClauses = append(setClauses, fmt.Sprintf("user_id = $%d", paramCounter))
		args = append(args, subscription.UserID)
		paramCounter++
	}

	if subscription.StartDate != "" {
		setClauses = append(setClauses, fmt.Sprintf("start_date = $%d", paramCounter))
		args = append(args, subscription.StartDate)
		paramCounter++
	}

	setClauses = append(setClauses, fmt.Sprintf("end_date = $%d", paramCounter))
	args = append(args, subscription.EndDate)
	paramCounter++

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", paramCounter))
	args = append(args, time.Now())
	paramCounter++

	args = append(args, id)

	query := fmt.Sprintf(
		"UPDATE subscriptions SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "),
		paramCounter,
	)

	result, err := r.db.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

// Delete deletes a subscription
func (r *SubscriptionRepository) Delete(id int) error {
	result, err := r.db.DB.Exec("DELETE FROM subscriptions WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

// CalculateTotalCost calculates the total cost of subscriptions for a period
func (r *SubscriptionRepository) CalculateTotalCost(req *models.CalculateCostRequest) (int, error) {
	query := `SELECT SUM(price) FROM subscriptions WHERE `
	whereConditions := []string{
		"(start_date <= $1 AND (end_date IS NULL OR end_date >= $2))",
	}
	args := []interface{}{req.EndPeriod, req.StartPeriod}
	paramCounter := 3

	if req.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("user_id = $%d", paramCounter))
		args = append(args, *req.UserID)
		paramCounter++
	}

	if req.ServiceName != nil && *req.ServiceName != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("service_name = $%d", paramCounter))
		args = append(args, *req.ServiceName)
		paramCounter++
	}

	query += strings.Join(whereConditions, " AND ")

	var totalCost sql.NullInt64
	err := r.db.DB.QueryRow(query, args...).Scan(&totalCost)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total cost: %w", err)
	}

	if !totalCost.Valid {
		return 0, nil
	}

	return int(totalCost.Int64), nil
}