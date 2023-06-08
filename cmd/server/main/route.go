package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func newRouter(app *App) http.Handler {
	router := httprouter.New()
	router.NotFound = app.notFound()

	router.HandlerFunc(http.MethodGet, "/home", app.home())
	router.HandlerFunc(http.MethodGet, "/signin", app.signinPage())
	router.HandlerFunc(http.MethodGet, "/login", app.signupPage())
	router.Handler(http.MethodGet, "/statics/*path", app.fileServer())

	staticMiddlewares := alice.New(app.recoverMiddleware)
	return staticMiddlewares.Then(router)
}
