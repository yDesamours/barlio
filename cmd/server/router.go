package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *App) routes() {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/home", Home(app))
	router.HandlerFunc(http.MethodGet, "/", NotFund(app))
	router.Handler(http.MethodGet, "/statics/", http.StripPrefix("/statics/", FileServer(app)))

	app.Server.Handler = recoverMiddleware(router)
}

func recoverMiddleware(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				http.Error(w, "an error has occured", http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	}
}
