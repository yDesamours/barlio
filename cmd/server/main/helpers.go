package main

import (
	"barlio/cmd/server/model"
	"net/http"
	"net/url"
)

func (app App) getUserHelper(r *http.Request) *model.User {
	user := r.Context().Value(userType).(*model.User)

	if user.ID == 0 {
		return nil
	}
	return user
}

func (app *App) readFormDataHelper(r *http.Request) url.Values {
	r.ParseForm()
	return r.Form
}

func (app *App) listArticleCategorieHelper() {}
