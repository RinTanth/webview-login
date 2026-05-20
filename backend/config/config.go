package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port           string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	FrontendOrigin string
	EmailPepper    string
	AESKey         string // 64-char hex string = 32 bytes for AES-256
	RedisHost             string
	RedisPort             string
	RedisPassword         string
	RefreshTokenTTL       time.Duration
	RedisRefreshTokenKey  string
}

func Load() *Config {
	return &Config{
		Port:           getenv("PORT", "8080"),
		DBHost:         mustEnv("DB_HOST"),
		DBPort:         mustEnv("DB_PORT"),
		DBUser:         mustEnv("DB_USER"),
		DBPassword:     mustEnv("DB_PASSWORD"),
		DBName:         mustEnv("DB_NAME"),
		JWTSecret:      mustEnv("JWT_SECRET"),
		FrontendOrigin: mustEnv("FRONTEND_ORIGIN"),
		EmailPepper:    mustEnv("EMAIL_PEPPER"),
		AESKey:         mustEnv("AES_KEY"),
		RedisHost:            mustEnv("REDIS_HOST"),
		RedisPort:            mustEnv("REDIS_PORT"),
		RedisPassword:        getenv("REDIS_PASSWORD", ""),
		RefreshTokenTTL:      mustDuration("REFRESH_TOKEN_TTL_DAYS") * 24 * time.Hour,
		RedisRefreshTokenKey: mustEnv("REDIS_REFRESH_TOKEN_KEY"),
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %q is not set", key)
	}
	return v
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustDuration(key string) time.Duration {
	v := mustEnv(key)
	days, err := strconv.Atoi(v)
	if err != nil || days <= 0 {
		log.Fatalf("environment variable %q must be a positive integer (days), got %q", key, v)
	}
	return time.Duration(days)
}
