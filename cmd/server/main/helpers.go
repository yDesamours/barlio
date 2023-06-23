package main

import (
	"barlio/cmd/server/model"
	"barlio/internal/token"
	"barlio/internal/types"
	"barlio/internal/validator"
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"time"
)

const (
	tokenDuration = time.Hour * 4
)

func (app App) getUser(r *http.Request) *model.User {
	user := r.Context().Value(userType).(*model.User)

	if user.ID == 0 {
		return nil
	}
	return user
}

func (app *App) ValidateUser(user *model.User, confirmedPassword string, validator *validator.Validator) {
	validator.NotEmpty(user.Username, "username", "missing username")
	if validator.NotEmpty(user.Email, "email", "missing email") {
		validator.IsEmailValid(user.Email, "email", "invalid email")
	}
	validator.NotEmpty(user.Password, "password", "missing password")
	if validator.NotEmpty(types.String(confirmedPassword), "passwordconfirm", "no password confirmation") {
		validator.Equal(user.Password, types.String(confirmedPassword), "password", "password mismatch")
	}
}

func (app *App) signInError(w http.ResponseWriter, data templateData, form url.Values, validator *validator.Validator) error {
	signinErrors := validator.Error()
	data.Set("errors", signinErrors)
	data.Set("signin", map[string]interface{}{"username": form.Get("username"), "email": form.Get("email")})
	data.Set("page", "Signin")

	tmpl := app.PageTemplates["signin"]
	return tmpl.Execute(w, data)
}

func (app *App) sendVerificationEmail(user *model.User, token *model.Token) {
	var (
		data         = templateData{"token": token.Token}
		mailTemplate = app.MailTemplate.Get("/emailverification")
	)

	mailObject, err := parseEmailData(mailTemplate, data)
	if err != nil {
		app.error(err)
		return
	}
	go app.Mailer.Send(string(user.Email), mailObject)
}

func (app *App) validateToken(userToken *model.Token, tokenString string, validator *validator.Validator) {
	validator.Check(userToken.ExpiretAt.After(time.Now()), "token", "token has expired")
	validator.Check(token.CompareToken(tokenString, userToken.Hash), "token", "token is invalid")
}

func parseEmailData(mailTemplate *PageTemplate, data templateData) (map[string]string, error) {
	var (
		buffer     = &bytes.Buffer{}
		mailObject = map[string]string{}
	)

	err := mailTemplate.ExecuteTemplate(buffer, "subject", data)
	if err != nil {
		return nil, err
	}
	mailObject["subject"] = buffer.String()
	buffer.Reset()

	err = mailTemplate.ExecuteTemplate(buffer, "alternative", data)
	if err != nil {
		return nil, err
	}
	mailObject["alternative"] = buffer.String()
	buffer.Reset()

	err = mailTemplate.ExecuteTemplate(buffer, "body", data)
	if err != nil {
		return nil, err
	}
	mailObject["body"] = buffer.String()

	return mailObject, nil
}

func (app *App) ValidateUserUnicity(u *model.User, validator *validator.Validator) error {
	user, err := app.models.user.Get(model.User{Username: types.String(u.Username)})
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return err
	}
	validator.Check(user.Username == "", "username", "username already taken")

	user, err = app.models.user.Get(model.User{Email: types.String(u.Email)})
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return err
	}
	validator.Check(user.Username == "", "email", "email already in use")
	return nil
}

func (app *App) newVerificationToken(user *model.User) (*model.Token, error) {
	plainTextToken, hashedToken, err := token.GenerateToken()
	if err != nil {
		return nil, err
	}
	token := model.Token{
		Userid:    user.ID,
		Scope:     model.EmailVerificationScope,
		Token:     plainTextToken,
		Hash:      hashedToken,
		ExpiretAt: time.Now().Add(tokenDuration),
	}
	return &token, nil
}

func (app *App) readFormData(r *http.Request) url.Values {
	r.ParseForm()
	return r.Form
}
