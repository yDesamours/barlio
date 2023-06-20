package main

import (
	"barlio/cmd/server/model"
	"context"
	"net/http"
)

func (app *App) recoverMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				app.panic(err)
				http.Error(w, "an error has occured", http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func (app *App) getCurrentUserMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := app.SessionManager.GetInt(r.Context(), "userId")
		user, err := app.models.user.Get(model.User{ID: userId})

		if err != nil {
			app.error(err)
			*user = model.NullUser()
		}

		ctx := context.WithValue(r.Context(), userType, user)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

func (app *App) notLoggedInOnlyMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.getUser(r)
		if user != nil && user.IsVerified {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (app *App) setNoCacheHeaderMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Expires", "0")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")

		h.ServeHTTP(w, r)
	})
}
