package app

import (
	"fmt"
	"net/http"
	"time"
)

func (app *App) handle() error {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(app.config.StaticPath))

	mux.HandleFunc("/", app.index)
	mux.HandleFunc("/search", app.search)
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	server := &http.Server{
		Addr:         ":" + app.config.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server.ListenAndServe()
}

func (app *App) index(w http.ResponseWriter, r *http.Request) {
	err := app.tmpl.ExecuteTemplate(w, "index", app.data.words)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *App) search(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Errorf("error parsing form: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	query := r.FormValue("q")
	if query == "" {
		app.tmpl.ExecuteTemplate(w, "results", app.data.words)
		return
	}
	if len(query) > 32 {
		query = query[:32]
	}

	if err := app.tmpl.ExecuteTemplate(w, "results", app.data.search(query)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
