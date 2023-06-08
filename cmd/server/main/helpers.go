package main

import (
	"barlio/internal/validator"
	"net/url"
)

func (app *App) validateSignFormHandler(form url.Values, validator *validator.Validator) {
	validator.NotEmpty(form.Get("username"), "username", "missing username")
	validator.NotEmpty(form.Get("email"), "email", "missing email")
	validator.IsEmailValid(form.Get("email"), "email", "invalid email")
	validator.NotEmpty(form.Get("password"), "password", "missing password")
	validator.NotEmpty(form.Get("sonfirmpassword"), "confirmpassword", "missing password confirmation")
	validator.Equal(form.Get("password"), form.Get("confirmpassword"), "confirmpassword", "must match password")
}
