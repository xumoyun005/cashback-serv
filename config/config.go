package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
	Env    string
}

type DBConfig struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     int
}

type ServerConfig struct {
	Port int
	Host string
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name,
	)
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	config := &Config{
		DB: DBConfig{
			Name:     getEnv("DB_NAME", "postgres"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
		},
		Server: ServerConfig{
			Port: serverPort,
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Env: getEnv("ENV", "development"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
