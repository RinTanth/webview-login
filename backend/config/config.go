package config

import (
	"log"
	"os"
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
