package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DBConfig
	Redis    RedisConfig
	S3       S3Config
	Auth     AuthConfig
}

type AppConfig struct {
	Port        string
	Environment string // development, production, test
	FrontendURL string
}

type DBConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

type S3Config struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
	PublicURL  string
	UseSSL     bool
}

type AuthConfig struct {
	JWTSecret string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found, using system environment variables")
	}

	return &Config{
		App: AppConfig{
			Port:        getEnv("PORT", "8080"),
			Environment: getEnv("APP_ENV", "development"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
		},
		Database: DBConfig{
			URL: getEnvOrFatal("DATABASE_URL"),
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6379/0"),
		},
		S3: S3Config{
			Endpoint:   getEnv("S3_ENDPOINT", ""),
			AccessKey:  getEnv("S3_ACCESS_KEY", ""),
			SecretKey:  getEnv("S3_SECRET_KEY", ""),
			BucketName: getEnv("S3_BUCKET_NAME", ""),
			PublicURL:  getEnv("S3_PUBLIC_URL", ""),
			UseSSL:     getEnvAsBool("S3_USE_SSL", false),
		},
		Auth: AuthConfig{
			JWTSecret: getEnvOrFatal("JWT_SECRET"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvOrFatal(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("FATAL: Environment variable %s is required but not set", key)
	}
	return value
}

func getEnvAsBool(key string, fallback bool) bool {
	valStr := getEnv(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return fallback
}
