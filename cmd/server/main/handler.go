package main

import (
	"barlio/cmd/server/model"
	"barlio/internal/data"
	"barlio/internal/mailer"
	"barlio/internal/validator"
	"barlio/ui"
	"net/http"
	"path/filepath"
	"text/template"
)

func (app *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Set("page", "Home")

	tmpl := app.Templates["home"]
	err := tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
	}
}

func (app *App) signinPageHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Set("page", "Signin")

	tmpl := app.Templates["signin"]
	err := tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
	}
}

func (app *App) signinHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.PostForm
	infos := app.newTemplateData()

	validator := validator.New()
	app.validateSignForm(form, validator)
	if !validator.Valid() {
		err := app.signInError(w, infos, form, validator)
		if err != nil {
			app.error(err)
		}
		return
	}

	user, err := app.User.Get(model.User{Username: data.String(form.Get("username"))})
	if err != nil {
		app.error(err)
		return
	}
	validator.Check(user == nil, "username", "username already taken")
	if !validator.Valid() {
		err := app.signInError(w, infos, form, validator)
		if err != nil {
			app.error(err)
		}
		return
	}

	user, err = app.User.Get(model.User{Email: data.String(form.Get("email"))})
	if err != nil {
		app.error(err)
		return
	}
	validator.Check(user == nil, "email", "email already in use")
	if !validator.Valid() {
		err := app.signInError(w, infos, form, validator)
		if err != nil {
			app.error(err)
		}
		return
	}

	user, err = app.newUser(form)
	if err != nil {
		app.error(err)
		return
	}

	err = app.User.Insert(user)
	if err != nil {
		app.error(err)
		return
	}

	go mailer.SendWelcomeEmail()

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (app *App) signupPageHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Set("page", "Signup")

	tmpl := app.Templates["signup"]
	err := tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
	}
}

func (app *App) notFound(w http.ResponseWriter, r *http.Request) {
	notFoundPath := filepath.Join("web", "html", "notFound.tmpl.html")
	temp, err := template.ParseFiles(notFoundPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	temp.Execute(w, nil)
}

func (app *App) fileServer() http.Handler {
	server := http.FileServer(http.FS(ui.FILES))
	return server
}
