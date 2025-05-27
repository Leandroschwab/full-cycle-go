package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	IPRateLimit    int
	TokenRateLimit int
	BlockDuration  int // em segundos
	RedisURL       string
}

func LoadConfig() *Config {
	// Carrega as variáveis de ambiente do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	config := &Config{}

	// Limite de requisições por IP
	ipLimit, err := strconv.Atoi(getEnv("IP_RATE_LIMIT", "10"))
	if err != nil {
		log.Fatalf("Invalid IP_RATE_LIMIT: %v", err)
	}
	config.IPRateLimit = ipLimit

	// Limite de requisições por token
	tokenLimit, err := strconv.Atoi(getEnv("TOKEN_RATE_LIMIT", "100"))
	if err != nil {
		log.Fatalf("Invalid TOKEN_RATE_LIMIT: %v", err)
	}
	config.TokenRateLimit = tokenLimit

	// Duração do bloqueio em segundos
	blockDuration, err := strconv.Atoi(getEnv("BLOCK_DURATION", "300"))
	if err != nil {
		log.Fatalf("Invalid BLOCK_DURATION: %v", err)
	}
	config.BlockDuration = blockDuration

	// URL do Redis
	config.RedisURL = getEnv("REDIS_URL", "localhost:6379")

	return config
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
