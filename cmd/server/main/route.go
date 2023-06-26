package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

const (
	HOMEPAGE              = "/"
	SIGNINPAGE            = "/signin"
	SIGNUPPAGE            = "/signup"
	EMAILVERIFICATIONPAGE = "/emailverification"
	SECURITYPAGE          = "/settings/security"
	EMAILCONFIRMROUTE     = "/emailconfirm"
	LOGOUTROUTE           = "/logout"
	ASSETSROUTE           = "/statics/*path"
	PROFILEPAGE           = "/settings/profile"
)

func newRouter(app *App) http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFound)

	notLoggedInOnly := alice.New(app.notLoggedInOnlyMiddleware)
	requireAuthMiddleware := alice.New(app.requireLoginMiddleware)

	router.HandlerFunc(http.MethodGet, HOMEPAGE, app.homePageHandler)
	router.Handler(http.MethodGet, SIGNINPAGE, notLoggedInOnly.ThenFunc(app.signinPageHandler))
	router.HandlerFunc(http.MethodGet, SIGNUPPAGE, app.signupPageHandler)
	router.Handler(http.MethodGet, EMAILVERIFICATIONPAGE, notLoggedInOnly.ThenFunc(app.emailVerificationHandler))
	router.Handler(http.MethodPost, SIGNINPAGE, notLoggedInOnly.ThenFunc(app.signinHandler))
	router.Handler(http.MethodPost, SIGNUPPAGE, notLoggedInOnly.ThenFunc(app.signupHandler))
	router.Handler(http.MethodPut, EMAILCONFIRMROUTE, notLoggedInOnly.ThenFunc(app.confirmEmailHandler))
	router.HandlerFunc(http.MethodGet, LOGOUTROUTE, app.logoutHandler)
	router.Handler(http.MethodGet, ASSETSROUTE, app.fileServer())
	router.Handler(http.MethodGet, PROFILEPAGE, requireAuthMiddleware.ThenFunc(app.profilePageHandler))

	staticMiddlewares := alice.New(app.recoverMiddleware, app.SessionManager.LoadAndSave, app.getCurrentUserMiddleware)
	staticMiddlewares.Append(app.setNoCacheHeaderMiddleware)
	return staticMiddlewares.Then(router)
}
