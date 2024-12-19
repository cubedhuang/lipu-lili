package app

import (
	"fmt"
	"html/template"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/cubedhuang/lipu-lili/internal/client"
	"github.com/cubedhuang/lipu-lili/internal/models"
)

type WordsStore struct {
	linkuData models.LinkuData
	words     []models.WordData
	mu        sync.RWMutex
}

type App struct {
	config *Config
	tmpl   *template.Template
	data   *WordsStore
}

func New(config *Config) (*App, error) {
	funcMap := template.FuncMap{
		"formatUsage": func(usage map[string]int) string {
			return fmt.Sprint(usage["2024-09"], "%")
		},
	}

	tmpl, err := template.New("base").Funcs(funcMap).ParseGlob(config.TemplatesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &App{
		config: config,
		tmpl:   tmpl,
		data:   &WordsStore{},
	}, nil
}

func (app *App) Run() error {
	if err := app.updateData(); err != nil {
		return fmt.Errorf("failed to fetch initial data: %w", err)
	}

	app.startUpdater()
	return app.handle()
}

func (app *App) updateData() error {
	newData, err := client.FetchLinku()
	if err != nil {
		return err
	}

	newWords := make([]models.WordData, 0, len(newData))
	for _, word := range newData {
		newWords = append(newWords, word)
	}

	slices.SortFunc(newWords, defaultSort)

	app.data.mu.Lock()
	defer app.data.mu.Unlock()

	app.data.linkuData = newData
	app.data.words = newWords

	return nil
}

func (app *App) startUpdater() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			if err := app.updateData(); err != nil {
				log.Printf("Failed to update data: %v", err)
			} else {
				log.Printf("Data updated at %v", time.Now())
			}
		}
	}()
}
