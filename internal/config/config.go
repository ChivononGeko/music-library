package config

import (
	"log"
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

func LoadConfig() *Config {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file, ensure it exists and contains the required values")
	}

	// Получаем переменные окружения
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatal("DB_HOST is not set in environment variables")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		log.Fatal("DB_PORT is not set in environment variables")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("DB_USER is not set in environment variables")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD is not set in environment variables")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME is not set in environment variables")
	}

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		log.Fatal("API_PORT is not set in environment variables")
	}

	externalAPI := os.Getenv("EXTERNAL_API_URL")
	if externalAPI == "" {
		log.Fatal("EXTERNAL_API_URL is not set in environment variables")
	}

	// Возвращаем конфигурацию
	return &Config{
		DBHost:      dbHost,
		DBPort:      dbPort,
		DBUser:      dbUser,
		DBPassword:  dbPassword,
		DBName:      dbName,
		APIPort:     apiPort,
		ExternalAPI: externalAPI,
	}
}
