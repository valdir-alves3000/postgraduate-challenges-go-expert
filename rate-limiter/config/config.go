package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort           string
	RedisHost         string
	RedisPort         string
	RequestLimitIP    int
	RequestLimitToken int
	BlockTime         int
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return Config{
		AppPort:           os.Getenv("APP_PORT"),
		RedisHost:         os.Getenv("REDIS_HOST"),
		RedisPort:         os.Getenv("REDIS_PORT"),
		RequestLimitIP:    getEnvAsInt("REQUEST_LIMIT_IP", 5),
		RequestLimitToken: getEnvAsInt("REQUEST_LIMIT_TOKEN", 10),
		BlockTime:         getEnvAsInt("BLOCK_TIME", 300),
	}
}

func getEnvAsInt(key string, defaultValue int) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return defaultValue
	}
	return val
}
