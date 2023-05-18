package main

import (
	"net/http"
	"path/filepath"
	"text/template"
)

func Home(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files := []string{
			filepath.Join("web", "html", "home.tmpl.html"),
			filepath.Join("web", "html", "header.tmpl.html"),
		}

		t, err := template.ParseFiles(files...)
		if err != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		err = t.ExecuteTemplate(w, "home", false)
		if err != nil {
			w.Write([]byte("welcome"))
			return
		}

	}
}

func NotFund(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notFoundPath := filepath.Join("web", "html", "notFound.tmpl.html")
		temp, err := template.ParseFiles(notFoundPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		temp.Execute(w, nil)
	}
}

func FileServer(app *App) http.Handler {
	server := http.FileServer(http.Dir(filepath.Join("web", "statics")))
	return server
}
