package main

import (
	"barlio/cmd/server/model"
	"barlio/internal/data"
	"barlio/internal/token"
	"barlio/internal/validator"
	"net/http"
	"net/url"
	"time"
)

const (
	tokenDuration = time.Hour * 4
)

func (app App) getUser(r *http.Request) *model.User {
	user := r.Context().Value("user").(*model.User)

	if user.ID == 0 {
		return nil
	}
	return user
}

func (app *App) newUser(form url.Values) *model.User {
	user := &model.User{
		Username: data.String(form.Get("username")),
		Email:    data.String(form.Get("email")),
		Password: data.String(form.Get("password")),
	}
	return user
}

func (app *App) ValidateUser(user *model.User, confirmedPassword string, validator *validator.Validator) {
	validator.NotEmpty(user.Username, "username", "missing username")
	if validator.NotEmpty(user.Email, "email", "missing email") {
		validator.IsEmailValid(user.Email, "email", "invalid email")
	}
	validator.NotEmpty(user.Password, "password", "missing password")
	if validator.NotEmpty(data.String(confirmedPassword), "passwordconfirm", "no password confirmation") {
		validator.Equal(user.Password, data.String(confirmedPassword), "password", "password mismatch")
	}
}

func (app *App) signInError(w http.ResponseWriter, data templateData, form url.Values, validator *validator.Validator) error {
	signinErrors := validator.Error()
	data.Set("errors", signinErrors)
	data.Set("signin", map[string]interface{}{"username": form.Get("username"), "email": form.Get("email")})
	data.Set("page", "Signin")

	tmpl := app.Templates["signin"]
	return tmpl.Execute(w, data)
}

func (app *App) sendWelcomeEmail(u *model.User) error {
	return nil
}

func (app *App) newVerificationToken(user *model.User) (*model.Token, error) {
	plainTextToken, hashedToken, err := token.GenerateToken()
	if err != nil {
		return nil, err
	}
	token := model.Token{
		Userid:    user.ID,
		Scope:     model.VerificationScope,
		Token:     plainTextToken,
		Hash:      hashedToken,
		ExpiretAt: time.Now().Add(tokenDuration),
	}
	return &token, nil
}
