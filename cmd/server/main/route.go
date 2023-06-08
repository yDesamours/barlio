package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func newRouter(app *App) http.Handler {
	router := httprouter.New()
	router.NotFound = app.notFound()

	router.HandlerFunc(http.MethodGet, "/home", app.homeHandler())
	router.HandlerFunc(http.MethodGet, "/signin", app.signinPageHandler())
	router.HandlerFunc(http.MethodGet, "/login", app.signupPageHandler())
	router.HandlerFunc(http.MethodGet, "/login", app.signinHandler())
	router.Handler(http.MethodGet, "/statics/*path", app.fileServer())

	staticMiddlewares := alice.New(app.recoverMiddleware)
	return staticMiddlewares.Then(router)
}
