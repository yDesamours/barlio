package main

import (
	"barlio/cmd/server/model"
	"context"
	"fmt"
	"net/http"
)

func (app App) recoverMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println(err)
				http.Error(w, "an error has occured", http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func (app App) getCurrentUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := app.SessionManager.GetInt(r.Context(), "userid")
		user, err := app.models.user.Get(model.User{ID: userId})

		if err != nil {
			app.error(err)
			*user = model.NullUser()
		}

		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})

}
