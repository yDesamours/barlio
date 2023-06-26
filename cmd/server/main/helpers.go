package main

import (
	"barlio/cmd/server/model"
	"barlio/internal/token"
	"barlio/internal/validator"
	"net/http"
	"net/url"
	"time"
)

func (app App) getUserHelper(r *http.Request) *model.User {
	return r.Context().Value(userType).(*model.User)
}

func (app *App) readFormDataHelper(r *http.Request) url.Values {
	r.ParseForm()
	return r.Form
}

func (app *App) validateTokenHelper(userToken *model.Token, tokenString string, validator *validator.Validator) {
	validator.Check(userToken.ExpiretAt.After(time.Now()), "token", "token has expired")
	validator.Check(token.CompareToken(tokenString, userToken.Hash), "token", "token is invalid")
}

func (app *App) listArticleCategorieHelper() {}
