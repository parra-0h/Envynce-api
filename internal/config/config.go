package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DSN       string
	AppEnv    string
	JWTSecret string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:      getEnv("PORT", "8080"),
		DSN:       getEnv("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=envynce_db port=5432 sslmode=disable"),
		AppEnv:    getEnv("APP_ENV", "development"),
		JWTSecret: getEnv("JWT_SECRET", "change-me-in-production"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
