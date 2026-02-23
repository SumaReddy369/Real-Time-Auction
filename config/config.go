package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort   string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	RedisAddr    string
	RedisDB      int
}

func Load() *Config {
	return &Config{
		ServerPort:   getEnv("PORT", "8080"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "auction"),
		DBPassword:   getEnv("DB_PASSWORD", "auction"),
		DBName:       getEnv("DB_NAME", "auction_db"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		RedisDB:      getEnvInt("REDIS_DB", 0),
	}
}

func (c *Config) DSN() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=disable"
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}
