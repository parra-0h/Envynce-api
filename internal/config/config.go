package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port   string
	DSN    string
	APIKey string
	AppEnv string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:   getEnv("PORT", "8080"),
		DSN:    getEnv("DATABASE_URL", "host=localhost user=postgres password=Hans989! dbname=config_service_db port=5432 sslmode=disable"),
		APIKey: getEnv("API_KEY", "secret-key"),
		AppEnv: getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
