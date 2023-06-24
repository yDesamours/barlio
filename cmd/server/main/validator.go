package main

import (
	"barlio/cmd/server/model"
	"barlio/internal/types"
	"barlio/internal/validator"
	"database/sql"
	"errors"
	"net/url"
	"time"
)

func (app *App) validateChangeUserPasswordForm(form url.Values, validator *validator.Validator) {
	validator.NotEmpty(types.String(form.Get("password")), "password", "password can't be empty")
	validator.NotEmpty(types.String(form.Get("passwordconfirm")), "passwordconfirm", "confirm your password")
	validator.Equal(types.String(form.Get("password")), types.String(form.Get("passwordconfirm")), "password", "password don't match")
}

func (app *App) validateUpdateUserProfileInfos(form url.Values, validator *validator.Validator) {
	_, err := time.Parse(time.DateOnly, form.Get("birthday"))
	validator.Check(err != nil, "birthday", "date is invalid")
}

func (app *App) validateUserUnicity(u *model.User, validator *validator.Validator) error {
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

func (app *App) validateUser(user *model.User, confirmedPassword string, validator *validator.Validator) {
	validator.NotEmpty(user.Username, "username", "missing username")
	if validator.NotEmpty(user.Email, "email", "missing email") {
		validator.IsEmailValid(user.Email, "email", "invalid email")
	}
	validator.NotEmpty(user.Password, "password", "missing password")
	if validator.NotEmpty(types.String(confirmedPassword), "passwordconfirm", "no password confirmation") {
		validator.Equal(user.Password, types.String(confirmedPassword), "password", "password mismatch")
	}
}
