package main

import (
	"barlio/cmd/server/model"
	"barlio/internal/helper"
	"barlio/internal/validator"
	"barlio/ui"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
)

func (app *App) homePageHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	user := app.getUser(r)

	data.Set("page", HOMEPAGE)
	data.Set("user", user)

	app.SessionManager.Put(r.Context(), "lastpage", HOMEPAGE)

	tmpl := app.PageTemplates.Get(HOMEPAGE)
	err := tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
	}
}

func (app *App) signinPageHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Set("showHeader", false)

	app.SessionManager.Put(r.Context(), "lastpage", SIGNINPAGE)

	tmpl := app.PageTemplates.Get(SIGNINPAGE)
	err := tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
	}
}

func (app *App) signinHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.PostForm
	infos := app.newTemplateData()

	user := app.models.user.NewUser(form)

	validator := validator.New()
	app.ValidateUser(user, form.Get("passwordconfirm"), validator)
	if !validator.Valid() {
		err := app.signInError(w, infos, form, validator)
		if err != nil {
			app.error(err)
		}
		return
	}

	err := app.ValidateUserUnicity(user, validator)
	if err != nil {
		app.error(err)
		return
	}

	hash, err := helper.HashPassword(user.Password)
	if err != nil {
		app.error(err)
		return
	}
	user.Password = hash

	err = app.models.user.Insert(user)
	if err != nil {
		app.error(err)
		return
	}

	app.SessionManager.Put(r.Context(), "userId", user.ID)

	http.Redirect(w, r, EMAILVERIFICATIONPAGE, http.StatusMovedPermanently)
}

func (app *App) emailVerificationHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Set("showHeader", false)

	user := app.getUser(r)
	if user == nil {
		http.Redirect(w, r, SIGNUPPAGE, http.StatusMovedPermanently)
		return
	}

	app.SessionManager.Put(r.Context(), "lastpage", EMAILVERIFICATIONPAGE)

	tmpl := app.PageTemplates.Get(EMAILVERIFICATIONPAGE)
	err := tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
		return
	}

	err = app.models.token.DeleteForUser(user.ID, model.EmailVerificationScope)
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

	go app.sendVerificationEmail(user, token)
}

func (app *App) signupPageHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Set("page", "Signup")
	data.Set("showHeader", false)

	app.SessionManager.Put(r.Context(), "lastpage", SIGNUPPAGE)

	tmpl := app.PageTemplates.Get(SIGNUPPAGE)
	err := tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
	}
}

func (app *App) signupHandler(w http.ResponseWriter, r *http.Request) {
	loginUser := app.models.user.NewUser(app.readFormData(r))
	templateData := templateData{"username": loginUser.Username, "password": ""}

	user, err := app.models.user.Get(*loginUser)
	if errors.Is(err, sql.ErrNoRows) {
		templateData["errors"] = map[string][]string{"username": {"username not found"}}
		tmpl := app.PageTemplates["signup"]
		err := tmpl.Execute(w, templateData)
		if err != nil {
			app.error(err)
		}
		return
	}

	if same := helper.CompareHash(loginUser.Password, user.Password); !same {
		templateData["errors"] = map[string][]string{"password": {"incorrect password"}}
		tmpl := app.PageTemplates["signup"]
		err := tmpl.Execute(w, templateData)
		if err != nil {
			app.error(err)
		}
		return
	}

	app.SessionManager.Put(r.Context(), "userId", user.ID)
	if !user.IsVerified {
		http.Redirect(w, r, EMAILVERIFICATIONPAGE, http.StatusMovedPermanently)
		return
	}

	http.Redirect(w, r, HOMEPAGE, http.StatusMovedPermanently)
}

func (app *App) confirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")

	if tokenString == "" {
		http.Redirect(w, r, EMAILVERIFICATIONPAGE, http.StatusMovedPermanently)
		return
	}

	user := app.getUser(r)
	if user == nil {
		http.Redirect(w, r, SIGNUPPAGE, http.StatusMovedPermanently)
		return
	}

	userToken, err := app.models.token.GetForUser(user.ID, model.EmailVerificationScope)
	if errors.Is(err, sql.ErrNoRows) {
		http.Redirect(w, r, EMAILVERIFICATIONPAGE, http.StatusMovedPermanently)
		return
	}

	validator := validator.New()
	app.validateToken(userToken, tokenString, validator)
	if !validator.Valid() {
		validator.Error()
		return
	}

	err = app.models.user.Activate(user)
	if err != nil {
		app.error(err)
		return
	}

	app.SessionManager.Put(r.Context(), "flash", "Account Verified")
	http.Redirect(w, r, HOMEPAGE, http.StatusMovedPermanently)
}

func (app *App) logoutHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Set("page", HOMEPAGE)

	err := app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.error(err)
		return
	}
	i := app.SessionManager.Pop(r.Context(), "userId")
	fmt.Println(i)
	http.Redirect(w, r, HOMEPAGE, http.StatusMovedPermanently)

}

func (app *App) notFound(w http.ResponseWriter, r *http.Request) {
	notFoundPath := filepath.Join("web", "html", "notFound.tmpl.html")
	temp, err := template.ParseFiles(notFoundPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	temp.Execute(w, templateData{})
}

func (app *App) fileServer() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = fmt.Sprintf("/web%s", r.URL.Path)
		server := http.FileServer(http.FS(ui.FILES))
		server.ServeHTTP(w, r)
	})

}
