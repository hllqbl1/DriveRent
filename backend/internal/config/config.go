package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func LoadConfig() (*Config, error) {
	// Membaca file .env.
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("gagal membaca file .env: %w", err)
	}

	// Mengambil konfigurasi dari environment.
	cfg := &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Validasi konfigurasi database.
	if cfg.DBHost == "" ||
		cfg.DBPort == "" ||
		cfg.DBUser == "" ||
		cfg.DBName == "" {
		return nil, fmt.Errorf("konfigurasi database belum lengkap")
	}

	// Validasi port database.
	port, err := strconv.Atoi(cfg.DBPort)
	if err != nil || port < 1 || port > 65535 {
		return nil, fmt.Errorf("DB_PORT tidak valid")
	}

	// Validasi port aplikasi.
	appPort, err := strconv.Atoi(cfg.AppPort)
	if err != nil || appPort < 1 || appPort > 65535 {
		return nil, fmt.Errorf("APP_PORT tidak valid")
	}

	return cfg, nil
}

// getEnv mengambil environment variable.
// Jika kosong, gunakan nilai default.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
