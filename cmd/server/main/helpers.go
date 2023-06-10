package main

import (
	"barlio/cmd/server/model"
	"barlio/internal/data"
	"barlio/internal/validator"
	"net/http"
	"net/url"
)

func (app *App) newUser(form url.Values) *model.User {
	user := &model.User{
		Username: data.String(form.Get("username")),
		Email:    data.String(form.Get("email")),
		Password: data.String(form.Get("password")),
	}
	return user
}

func (app *App) ValidateUser(user *model.User, validator *validator.Validator) {
	validator.NotEmpty(user.Username, "username", "missing username")
	if validator.NotEmpty(user.Email, "email", "missing email") {
		validator.IsEmailValid(user.Email, "email", "invalid email")
	}
	validator.NotEmpty(user.Password, "password", "missing password")
}

func (app *App) signInError(w http.ResponseWriter, data templateData, form url.Values, validator *validator.Validator) error {
	signinErrors := validator.Error()
	data.Set("errors", signinErrors)
	data.Set("signin", map[string]interface{}{"username": form.Get("username"), "email": form.Get("email")})
	data.Set("page", "Signin")

	tmpl := app.Templates["signin"]
	return tmpl.Execute(w, data)
}
