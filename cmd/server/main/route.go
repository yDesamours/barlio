package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func newRouter(app *App) http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFound)

	notLoggedInOnly := alice.New(app.notLoggedInOnly)

	router.HandlerFunc(http.MethodGet, "/", app.homeHandler)
	router.HandlerFunc(http.MethodGet, "/signin", app.signinPageHandler)
	router.Handler(http.MethodGet, "/login", notLoggedInOnly.ThenFunc(app.signupPageHandler))
	router.Handler(http.MethodGet, "/emailverification", notLoggedInOnly.ThenFunc(app.emailVerificationPageHandler))
	router.HandlerFunc(http.MethodPost, "/signin", app.signinHandler)
	router.Handler(http.MethodGet, "/statics/*path", app.fileServer())

	staticMiddlewares := alice.New(app.recoverMiddleware, app.SessionManager.LoadAndSave, app.getCurrentUserMiddleware)
	return staticMiddlewares.Then(router)
}
