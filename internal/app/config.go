package app

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Port           string
	UpdateInterval time.Duration
	TemplatesPath  string
	StaticPath     string
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	interval := os.Getenv("UPDATE_INTERVAL")
	if interval == "" {
		interval = "10m"
	}

	updateInterval, err := time.ParseDuration(interval)
	if err != nil {
		return nil, fmt.Errorf("failed to parse update interval: %w", err)
	}

	return &Config{
		Port:           port,
		UpdateInterval: updateInterval,
		TemplatesPath:  "templates/*.html",
		StaticPath:     "static",
	}, nil
}
