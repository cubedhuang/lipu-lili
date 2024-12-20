package app

import (
	"net/http"

	"github.com/cubedhuang/lipu-lili/internal/models"
)

func (app *App) createHandlers() *http.ServeMux {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(app.config.StaticPath))

	mux.HandleFunc("/", app.index)
	mux.HandleFunc("/{id}", app.word)
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux
}

type IndexPage struct {
	Query   string
	Results []models.WordData
}

func (app *App) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	query := ""
	if err := r.ParseForm(); err == nil {
		query = r.FormValue("q")
	}

	results := app.data.search(query)

	if r.Header.Get("Hx-Request") == "true" {
		if err := app.tmpl.ExecuteTemplate(w, "results", results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := app.tmpl.ExecuteTemplate(w, "index", IndexPage{
		Query:   query,
		Results: results,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type WordPage struct {
	Word  models.WordData
	Signs []models.SignData
}

func (app *App) word(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	word, ok := app.data.linkuData.Words[id]
	if !ok {
		http.NotFound(w, r)
		return
	}

	if err := app.tmpl.ExecuteTemplate(w, "word", WordPage{
		Word:  word,
		Signs: app.data.signs[id],
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
