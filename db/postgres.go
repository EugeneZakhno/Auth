package db

import (
	"database/sql"
	"fmt"
	"subscription-service/config"

	_ "github.com/lib/pq"
)

// PostgresDB represents a PostgreSQL database connection
type PostgresDB struct {
	DB *sql.DB
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(cfg config.DatabaseConfig) (*PostgresDB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{DB: db}, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	return p.DB.Close()
}

// RunMigrations runs database migrations
func (p *PostgresDB) RunMigrations() error {
	// Create subscriptions table
	_, err := p.DB.Exec(`
		CREATE TABLE IF NOT EXISTS subscriptions (
			id SERIAL PRIMARY KEY,
			service_name VARCHAR(255) NOT NULL,
			price INTEGER NOT NULL,
			user_id UUID NOT NULL,
			start_date VARCHAR(7) NOT NULL,
			end_date VARCHAR(7),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create subscriptions table: %w", err)
	}

	// Insert some test data if the table is empty
	var count int
	err = p.DB.QueryRow("SELECT COUNT(*) FROM subscriptions").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count subscriptions: %w", err)
	}

	if count == 0 {
		_, err = p.DB.Exec(`
			INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
			VALUES 
			('Netflix', 599, '60601fee-2bf1-4721-ae6f-7636e79a0cba', '01-2023', '01-2024'),
			('Spotify', 199, '60601fee-2bf1-4721-ae6f-7636e79a0cba', '02-2023', NULL),
			('Yandex Plus', 299, '70701fee-3bf1-5721-be6f-8636e79a0cba', '03-2023', '03-2024')
		`)
		if err != nil {
			return fmt.Errorf("failed to insert test data: %w", err)
		}
	}

	return nil
}