package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	APIPort     string
	ExternalAPI string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file, ensure it exists and contains the required values")
		return nil, err
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		return nil, fmt.Errorf("the DB_HOST value is not set in the environment variables")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		return nil, fmt.Errorf("the DB_PORT value is not set in the environment variables")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return nil, fmt.Errorf("the DB_USER value is not set in the environment variables")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	// if dbPassword == "" {
	// 	return nil, fmt.Errorf("the DB_PASSWORD value is not set in the environment variables")
	// }

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, fmt.Errorf("the DB_NAME value is not set in the environment variables")
	}

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		return nil, fmt.Errorf("the API_PORT value is not set in the environment variables")
	}

	externalAPI := os.Getenv("EXTERNAL_API_URL")
	// if externalAPI == "" {
	// 	return nil, fmt.Errorf("the EXTERNAL_API_URL value is not set in the environment variables")
	// }

	return &Config{
		DBHost:      dbHost,
		DBPort:      dbPort,
		DBUser:      dbUser,
		DBPassword:  dbPassword,
		DBName:      dbName,
		APIPort:     apiPort,
		ExternalAPI: externalAPI,
	}, nil
}
