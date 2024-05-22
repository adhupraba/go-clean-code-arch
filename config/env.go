package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost string
	ServerPort string

	DbUrl string

	JWTSecret              string
	JWTExpirationInSeconds int64
}

var Envs = initConfig()

func initConfig() Config {
	_, b, _, _ := runtime.Caller(0)
	envPath := filepath.Join(filepath.Dir(b), "../.env")

	err := godotenv.Load(envPath)

	if err != nil {
		log.Fatal("Unable to load .env:", err)
	}

	log.Println("env loaded")

	dbUrl := getEnv("DB_URL", "")

	if dbUrl == "" {
		log.Fatal("db url not present")
	}

	return Config{
		ServerHost: getEnv("SERVER_HOST", "http://localhost"),
		ServerPort: getEnv("SERVER_PORT", "8080"),

		DbUrl: dbUrl,

		JWTSecret:              getEnv("JWT_SECRET", "secret"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXP", 3600*24*7),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
