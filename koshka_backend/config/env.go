package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// AppConfig stores global configuration
type AppConfig struct {
	BaseDomain string
}

// Config holds the loaded configuration
var Config AppConfig

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	// Assign environment variables to config
	Config = AppConfig{
		BaseDomain: os.Getenv("BASE_DOMAIN"),
	}

	// Ensure BaseDomain is set
	if Config.BaseDomain == "" {
		log.Fatal("Error: BASE_DOMAIN is not set")
	}
}

// ✅ Get JWT Secret
func GetJWTSecret(Secret string) string {
	return os.Getenv(Secret)
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
