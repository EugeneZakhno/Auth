package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

// LoadConfig loads the application configuration from environment variables and config file
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Default configuration
	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "subscriptions"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	// Try to load config from YAML file
	configPath := getEnv("CONFIG_PATH", "config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		if err := loadYAMLConfig(configPath, cfg); err != nil {
			return nil, fmt.Errorf("failed to load YAML config: %w", err)
		}
	}

	return cfg, nil
}

// loadYAMLConfig loads configuration from a YAML file
func loadYAMLConfig(path string, cfg *Config) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}