package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

var Module = fx.Module("config", fx.Provide(Load))

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Upload   UploadConfig
	Markdown MarkdownConfig
}

type ServerConfig struct {
	Addr string
}

type DatabaseConfig struct {
	DSN string
}

type UploadConfig struct {
	SignatureSecret  string
	MaxMarkdownBytes int64
	TimestampSkew    time.Duration
	AllowedIPs       []string
}

type MarkdownConfig struct {
	RenderHTML   bool
	SanitizeHTML bool
}

func Load() (Config, error) {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("backend/.env")

	cfg := Config{
		Server: ServerConfig{
			Addr: getString("SERVER_ADDR", "127.0.0.1:8080"),
		},
		Database: DatabaseConfig{
			DSN: getString("DATABASE_DSN", ""),
		},
		Upload: UploadConfig{
			SignatureSecret:  getString("UPLOAD_SIGNATURE_SECRET", ""),
			MaxMarkdownBytes: 1024 * 1024,
			TimestampSkew:    5 * time.Minute,
		},
		Markdown: MarkdownConfig{
			RenderHTML:   true,
			SanitizeHTML: true,
		},
	}

	if cfg.Database.DSN == "" {
		return Config{}, fmt.Errorf("DATABASE_DSN is required")
	}
	if cfg.Upload.SignatureSecret == "" {
		return Config{}, fmt.Errorf("UPLOAD_SIGNATURE_SECRET is required")
	}

	return cfg, nil
}

func getString(key string, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
