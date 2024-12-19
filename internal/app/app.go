package app

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/cubedhuang/lipu-lili/internal/client"
	"github.com/cubedhuang/lipu-lili/internal/models"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
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

	tmpl, err := compileTemplates(config.TemplatesPath, funcMap)
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

func compileTemplates(path string, funcMap template.FuncMap) (*template.Template, error) {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)

	tmpl := template.New("").Funcs(funcMap)

	files, err := filepath.Glob(path)
	if err != nil {
		return nil, fmt.Errorf("failed to find templates: %w", err)
	}

	for _, filename := range files {
		b, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
		}

		mb, err := m.Bytes("text/html", b)
		if err != nil {
			return nil, fmt.Errorf("failed to minify file %s: %w", filename, err)
		}

		tmpl, err = tmpl.Parse(string(mb))
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", filename, err)
		}
	}

	return tmpl, nil
}
