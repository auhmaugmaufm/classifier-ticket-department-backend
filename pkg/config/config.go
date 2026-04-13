package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string

	JWTSecret     string
	JWTExpireHour int

	AIBackendUrl string
	// AIApiKey     string
}

var cfg *Config

func Load() {
	if os.Getenv("APP_NEW") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("no .env file found")
		}
	}

	expireHour, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOUR"))
	cfg = &Config{
		AppEnv:  os.Getenv("APP_ENV"),
		AppPort: os.Getenv("APP_PORT"),

		DBHost:    os.Getenv("DB_HOST"),
		DBPort:    os.Getenv("DB_PORT"),
		DBUser:    os.Getenv("DB_USER"),
		DBPass:    os.Getenv("DB_PASS"),
		DBName:    os.Getenv("DB_NAME"),
		DBSSLMode: os.Getenv("DB_SSLMODE"),

		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpireHour: expireHour,

		AIBackendUrl: os.Getenv("AI_BASE_URL"),
		// AIApiKey:     os.Getenv("Ai_API_KEY"),
	}
}

func Get() *Config {
	return cfg
}
