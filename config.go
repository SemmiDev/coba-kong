package main

import (
	"os"
)

// Config menyimpan semua konfigurasi aplikasi yang dibaca dari environment variables
// Menggunakan struct ini membuat kode lebih maintainable karena semua config terpusat
type Config struct {
	ServiceName  string
	Environment  string
	Port         string
	DatabaseURL  string
	JWTSecret    string
	KongAdminURL string
}

// LoadConfig membaca environment variables dan mengembalikan Config struct
// Ini menggunakan pattern yang umum di Go untuk centralized configuration
func LoadConfig() *Config {
	return &Config{
		ServiceName:  getEnv("SERVICE_NAME", "user-service"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		KongAdminURL: getEnv("KONG_ADMIN_URL", "http://kong:8001"),
	}
}

// getEnv adalah helper function untuk membaca env var dengan default value
// Ini mencegah aplikasi crash jika env var tidak di-set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
