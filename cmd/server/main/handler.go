package main

import (
	"barlio/ui"
	"net/http"
	"path/filepath"
	"text/template"
)

func (app *App) home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := newTemplateData()
		data.Set("page", "Home")
		tmpl := app.Templates["home"]
		err := tmpl.Execute(w, data)
		if err != nil {
			app.error(err)
		}
	}
}

func (app *App) signinPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := newTemplateData()
		data.Set("page", "Signin")
		tmpl := app.Templates["signin"]
		err := tmpl.Execute(w, data)
		if err != nil {
			app.error(err)
		}
	}
}

func (app *App) notFound() http.HandlerFunc {
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

func (app *App) fileServer() http.Handler {
	server := http.FileServer(http.FS(ui.FILES))
	return server
}
