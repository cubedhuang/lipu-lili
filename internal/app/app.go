package app

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/cubedhuang/lipu-lili/internal/client"
	"github.com/cubedhuang/lipu-lili/internal/models"
)

type WordsStore struct {
	linkuData *models.LinkuData
	words     []models.WordData
	signs     map[string][]models.SignData
	mu        sync.RWMutex
}

type App struct {
	config *Config
	tmpl   *template.Template
	data   *WordsStore
}

func New(config *Config) (*App, error) {
	funcMap := template.FuncMap{
		"join": strings.Join,
		"formatUsage": func(usage map[string]int) string {
			val, ok := usage["2024-09"]
			if !ok {
				return "unknown"
			}
			return fmt.Sprint(val, "%")
		},
		"processPuData": func(data models.WordPuVerbatim) []models.PuData {
			puData := make([]models.PuData, 0)

			lines := strings.Split(data.En, "\n")
			for _, line := range lines {
				parts := strings.SplitN(line, " ", 2)
				puData = append(puData, models.PuData{
					PartOfSpeech: parts[0],
					Definition:   parts[1],
				})
			}

			return puData
		},
		"processEtymologyData": func(word models.WordData) models.EtymologyData {
			data := models.EtymologyData{
				Source:  "",
				Entries: make([]models.EtymologyEntry, 0, len(word.Etymology)),
			}

			if strings.HasPrefix(word.SourceLanguage, "multiple") || strings.HasPrefix(word.SourceLanguage, "unknown") {
				data.Source = word.SourceLanguage
			}

			for i, entry := range word.Etymology {
				data.Entries = append(data.Entries, models.EtymologyEntry{
					Word:       entry.Word,
					Alt:        entry.Alt,
					Definition: word.Translations["en"].Etymology[i].Definition,
					Language:   word.Translations["en"].Etymology[i].Language,
				})
			}

			return data
		},
		"fromCodePoint": func(ucsur string) string {
			// e.g. U+F1900
			codePoint := strings.TrimPrefix(ucsur, "U+")
			code, err := strconv.ParseInt(codePoint, 16, 32)
			if err != nil {
				return ""
			}
			return string(rune(code))
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app.startUpdater(ctx)
	mux := app.createHandlers()

	server := &http.Server{
		Addr:         ":" + app.config.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	errors := make(chan error, 1)

	go func() {
		errors <- server.ListenAndServe()
	}()

	select {
	case err := <-errors:
		return fmt.Errorf("server error: %w", err)
	case <-shutdown:
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}

		log.Println("Server gracefully stopped")
		return nil
	}
}

func (app *App) updateData() error {
	newData, err := client.FetchLinku()
	if err != nil {
		return err
	}

	newWords := make([]models.WordData, 0, len(newData.Words))
	for _, word := range newData.Words {
		newWords = append(newWords, word)
	}

	slices.SortFunc(newWords, defaultSort)

	newSigns := make(map[string][]models.SignData, len(newWords))
	for _, word := range newWords {
		for _, sign := range newData.Signs {
			if sign.Definition != word.Id {
				continue
			}

			_, ok := newSigns[word.Id]
			if !ok {
				newSigns[word.Id] = make([]models.SignData, 0)
			}

			newSigns[word.Id] = append(newSigns[word.Id], sign)
		}

		slices.SortFunc(newSigns[word.Id], func(a, b models.SignData) int {
			return strings.Compare(a.Id, b.Id)
		})
	}

	app.data.mu.Lock()
	defer app.data.mu.Unlock()

	app.data.linkuData = newData
	app.data.words = newWords
	app.data.signs = newSigns

	return nil
}

func (app *App) startUpdater(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := app.updateData(); err != nil {
					log.Printf("Failed to update data: %v", err)
				}
			}
		}
	}()
}
