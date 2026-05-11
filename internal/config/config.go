package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL     string
	PORT      string
	JWTSECRET string
}

func LoadEnv() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Info: missing .env file")
	}
	return &Config{
		DBURL:     getEnv("DATABASE_URL", ""),
		PORT:      getEnv("PORT", "8080"),
		JWTSECRET: getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
