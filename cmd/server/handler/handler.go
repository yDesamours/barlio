package handler

import (
	"barlio/cmd/server/config"
	"net/http"
	"path/filepath"
	"text/template"
)

func Home(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	}
}

func NotFund(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notFoundPath := filepath.Join(app.Path, "web", "html", "notFound.tmpl.html")
		temp, err := template.ParseFiles(notFoundPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		temp.Execute(w, nil)
	}
}
