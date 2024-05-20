package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost string
	ServerPort string

	DbUrl string
}

var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Unable to load .env:", err)
	}

	dbUrl := getEnv("DB_URL", "")

	if dbUrl == "" {
		log.Fatal("db url not present")
	}

	return Config{
		ServerHost: getEnv("SERVER_HOST", "http://localhost"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DbUrl:      dbUrl,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
