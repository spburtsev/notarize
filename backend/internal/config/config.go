package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	LogLevel    string
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3UseSSL    bool
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
		S3Endpoint:  os.Getenv("S3_ENDPOINT"),
		S3AccessKey: os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey: os.Getenv("S3_SECRET_KEY"),
		S3Bucket:    os.Getenv("S3_BUCKET"),
		S3UseSSL:    os.Getenv("S3_USE_SSL") == "true",
	}
	if cfg.DatabaseURL == "" {
		return cfg, errors.New("DATABASE_URL is required")
	}
	if cfg.S3Endpoint == "" {
		return cfg, errors.New("S3_ENDPOINT is required")
	}
	if cfg.S3AccessKey == "" {
		return cfg, errors.New("S3_ACCESS_KEY is required")
	}
	if cfg.S3SecretKey == "" {
		return cfg, errors.New("S3_SECRET_KEY is required")
	}
	if cfg.S3Bucket == "" {
		return cfg, errors.New("S3_BUCKET is required")
	}
	return cfg, nil
}
