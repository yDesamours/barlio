package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func newRouter(app *App) http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFound)

	router.HandlerFunc(http.MethodGet, "/", app.homeHandler)
	router.HandlerFunc(http.MethodGet, "/signin", app.signinPageHandler)
	router.HandlerFunc(http.MethodGet, "/login", app.signupPageHandler)
	router.HandlerFunc(http.MethodGet, "/verification", app.verificationHandler)
	router.HandlerFunc(http.MethodPost, "/signin", app.signinHandler)
	router.Handler(http.MethodGet, "/statics/*path", app.fileServer())

	staticMiddlewares := alice.New(app.recoverMiddleware, app.SessionManager.LoadAndSave, app.getCurrentUserMiddleware)
	return staticMiddlewares.Then(router)
}
