package main

import (
	"barlio/internal/validator"
	"barlio/ui"
	"net/http"
	"path/filepath"
	"text/template"
)

func (app *App) homeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.newTemplateData()
		data.Set("page", "Home")

		tmpl := app.Templates["home"]
		err := tmpl.Execute(w, data)
		if err != nil {
			app.error(err)
		}
	}
}

func (app *App) signinPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.newTemplateData()
		data.Set("page", "Signin")

		tmpl := app.Templates["signin"]
		err := tmpl.Execute(w, data)
		if err != nil {
			app.error(err)
		}
	}
}

func (app *App) signinHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		form := r.PostForm
		data := app.newTemplateData()

		validator := &validator.Validator{}
		app.validateSignFormHandler(form, validator)
		if !validator.Valid() {
			signinErrors := validator.Error()
			data.Set("errors", signinErrors)
			data.Set("signin", map[string]interface{}{"username": form.Get("username"), "email": form.Get("email")})
			data.Set("page", "Signin")

			tmpl := app.Templates["signin"]
			err := tmpl.Execute(w, data)
			if err != nil {
				app.error(err)
			}

		}

	}
}

func (app *App) signupPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.newTemplateData()
		data.Set("page", "Signup")

		tmpl := app.Templates["signup"]
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
