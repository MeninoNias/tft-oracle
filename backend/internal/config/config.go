package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
	ServerHost  string
	AppEnv      string
	LogLevel    string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		ServerHost:  getEnv("SERVER_HOST", "localhost"),
		AppEnv:      getEnv("APP_ENV", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func (c *Config) ListenAddr() string {
	return fmt.Sprintf("%s:%s", c.ServerHost, c.ServerPort)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
