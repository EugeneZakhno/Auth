# Subscription Management Service

A REST service for aggregating data on users' online subscriptions.

## Features

- CRUD operations for subscription records
- Calculate total cost of subscriptions for a selected period
- Filter subscriptions by user ID and service name
- PostgreSQL database with migrations
- Swagger documentation
- Docker deployment

## Tech Stack

- Go 1.21
- Gin Web Framework
- PostgreSQL
- Docker & Docker Compose

## API Endpoints

- `POST /api/v1/subscriptions` - Create a new subscription
- `GET /api/v1/subscriptions` - List all subscriptions
- `GET /api/v1/subscriptions/:id` - Get a subscription by ID
- `PUT /api/v1/subscriptions/:id` - Update a subscription
- `DELETE /api/v1/subscriptions/:id` - Delete a subscription
- `GET /api/v1/subscriptions/calculate` - Calculate total subscription cost

## Running the Application

### Prerequisites

- Docker and Docker Compose installed

### Using Docker Compose

1. Clone the repository

2. Start the application and database:

```bash
docker-compose up -d
```

3. The service will be available at http://localhost:8080

4. Access Swagger documentation at http://localhost:8080/swagger/index.html

### Development Setup

1. Install Go 1.21 or later

2. Install PostgreSQL

3. Create a database named `subscriptions`

4. Set up environment variables or modify `config.yaml`

5. Install dependencies:

```bash
go mod download
```

6. Generate Swagger documentation:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

7. Run the application:

```bash
go run main.go
```

## Configuration

The application can be configured using environment variables or a `config.yaml` file. Environment variables take precedence over the config file.

### Environment Variables

- `SERVER_HOST` - Server host (default: "localhost")
- `SERVER_PORT` - Server port (default: "8080")
- `DB_HOST` - Database host (default: "localhost")
- `DB_PORT` - Database port (default: "5432")
- `DB_USER` - Database user (default: "postgres")
- `DB_PASSWORD` - Database password (default: "postgres")
- `DB_NAME` - Database name (default: "subscriptions")
- `DB_SSLMODE` - Database SSL mode (default: "disable")
- `LOG_LEVEL` - Logging level (default: "info")

## Example Requests

### Create a Subscription

```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

### Calculate Total Cost

```bash
curl -X GET "http://localhost:8080/api/v1/subscriptions/calculate?start_period=01-2023&end_period=12-2023&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba"
```

## Database

The application automatically runs migrations on startup to create the necessary tables and inserts test data if the tables are empty.