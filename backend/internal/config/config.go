package config

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	MediaDir    string
}

func Load() Config {
	return Config{
		Port: getEnv("PORT", "8080"),
		// DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5430/localconnect?sslmode=disable"),
		DatabaseURL: getEnv("DATABASE_URL", "postgresql://postgres:RishuRaj%402002@db.ggtbvdxbwuokfebhwwyq.supabase.co:5432/postgres"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change"),
		MediaDir:    getEnv("MEDIA_DIR", "./media"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
