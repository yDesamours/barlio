package main

import (
	"barlio/cmd/server/model"
	"barlio/internal/token"
	"barlio/internal/validator"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

func (app App) getUserHelper(r *http.Request) *model.User {
	return r.Context().Value(userType).(*model.User)
}

func (app *App) setLastPageHelper(r *http.Request) {
	app.SessionManager.Put(r.Context(), "lastpage", r.URL.Path)
}

func (app *App) readMultipartFormDataHelper(r *http.Request) (*multipart.Form, error) {
	maxBytes := 8 << 20
	err := r.ParseMultipartForm(int64(maxBytes))
	if err != nil {
		return nil, err
	}
	return r.MultipartForm, nil
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
