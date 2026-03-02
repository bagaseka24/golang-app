package internal

import (
	"os"
)

type Config struct {
	AppPort     string
	DbEnabled   bool
	DbURL       string
	S3Enabled   bool
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
}

func LoadConfig() Config {
	return Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		DbEnabled:   getEnv("DB_ENABLED", "false") == "true",
		DbURL:       getEnv("DATABASE_URL", ""),
		S3Enabled:   getEnv("S3_ENABLED", "false") == "true",
		S3Endpoint:  getEnv("S3_ENDPOINT", ""),
		S3AccessKey: getEnv("S3_ACCESS_KEY", ""),
		S3SecretKey: getEnv("S3_SECRET_KEY", ""),
		S3Bucket:    getEnv("S3_BUCKET", "uploads"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
