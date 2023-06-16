package main

import (
	"barlio/cmd/server/model"
	"barlio/internal/validator"
	"barlio/ui"
	"net/http"
	"path/filepath"
	"text/template"
)

func (app *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	user := app.getUser(r)

	data.Set("page", "Home")
	data.Set("user", user)

	tmpl := app.Templates["home"]
	err := tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
	}
}

func (app *App) signinPageHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Set("page", "Signin")
	data.Set("showHeader", false)

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

	user := app.newUser(form)

	validator := validator.New()
	app.ValidateUser(user, form.Get("passwordconfirm"), validator)
	if !validator.Valid() {
		err := app.signInError(w, infos, form, validator)
		if err != nil {
			app.error(err)
		}
		return
	}

	err := app.models.user.ValidateUser(user, validator)
	if err != nil {
		app.error(err)
		return
	}

	err = user.HashPassword()
	if err != nil {
		app.error(err)
		return
	}

	err = app.models.user.Insert(user)
	if err != nil {
		app.error(err)
		return
	}

	app.SessionManager.Put(r.Context(), "userId", user.ID)

	http.Redirect(w, r, "/verification", http.StatusMovedPermanently)
}

func (app *App) emailVerificationPageHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Set("page", "emailverification")
	data.Set("showHeader", false)

	user := app.getUser(r)
	if user == nil {
		http.Redirect(w, r, "/signup", http.StatusMovedPermanently)
		return
	}

	if user.IsVerified {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	tmpl := app.Templates["emailverification"]
	err := tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
	}
}

func (app *App) emailVerificationTokenPageHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Set("page", "emailverificationtoken")
	data.Set("showHeader", false)

	user := app.getUser(r)
	if user == nil {
		http.Redirect(w, r, "/signup", http.StatusMovedPermanently)
		return
	}

	if user.IsVerified {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	tmpl := app.Templates["emailverificationtoken"]
	err := tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
		return
	}

	err = app.models.token.DeleteForUser(user.ID, model.VerificationScope)
	if err != nil {
		app.error(err)
		return
	}

	token, err := app.newVerificationToken(user)
	if err != nil {
		app.error(err)
		return
	}

	err = app.models.token.Insert(token)
	if err != nil {
		app.error(err)
		return
	}
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
