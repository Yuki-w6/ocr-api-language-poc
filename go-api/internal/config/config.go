package config

import "os"

type Config struct {
	Port                  string
	DatabaseURL           string
	UploadURLBase         string
	PresignedURLExpiresIn int
}

func Load() Config {
	return Config{
		Port:                  getEnv("PORT", "8080"),
		DatabaseURL:           getEnv("DATABASE_URL", "postgres://postgres:postgres@db:5432/ocr_poc?sslmode=disable"),
		UploadURLBase:         getEnv("UPLOAD_URL_BASE", "https://storage.example.com/upload"),
		PresignedURLExpiresIn: 300,
	}
}

func getEnv(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}
