package main

import (
	"barlio/cmd/server/model"
	"barlio/internal/data"
	"barlio/internal/validator"
	"net/http"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

func (app *App) validateSignForm(form url.Values, validator *validator.Validator) {
	validator.NotEmpty(form.Get("username"), "username", "missing username")
	if validator.NotEmpty(form.Get("email"), "email", "missing email") {
		validator.IsEmailValid(form.Get("email"), "email", "invalid email")
	}
	validator.NotEmpty(form.Get("password"), "password", "missing password")
	if validator.NotEmpty(form.Get("passwordconfirm"), "passwordconfirm", "missing password confirmation") {
		validator.Equal(form.Get("password"), form.Get("passwordconfirm"), "passwordconfirm", "must match password")
	}
}

func (app *App) newUser(form url.Values) (*model.User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(form.Get("password")), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Username: data.String(form.Get("username")),
		Email:    data.String(form.Get("email")),
		Password: data.String(hashPassword),
	}
	return user, nil
}

func (app *App) signInError(w http.ResponseWriter, data templateData, form url.Values, validator *validator.Validator) error {
	signinErrors := validator.Error()
	data.Set("errors", signinErrors)
	data.Set("signin", map[string]interface{}{"username": form.Get("username"), "email": form.Get("email")})
	data.Set("page", "Signin")

	tmpl := app.Templates["signin"]
	return tmpl.Execute(w, data)
}
