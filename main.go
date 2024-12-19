package main

import (
	"log"

	"github.com/cubedhuang/lipu-lili/internal/app"
)

func main() {
	config, err := app.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	application, err := app.New(config)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("application error: %v", err)
	}
}
