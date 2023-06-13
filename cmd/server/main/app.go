package main

import (
	"barlio/cmd/server/model"
	"log"

	"database/sql"
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type App struct {
	Server         *http.Server
	SessionManager *scs.SessionManager
	DB             *sql.DB
	Logger         *Logger
	Config         *Config
	models         *models
	Templates      map[string]*PageTemplate
}

type models struct {
	user    *model.UserModel
	article *model.ArticleModel
	token   *model.TokenModel
}

type Config struct {
	Server struct {
		Port int
	}
	DB struct {
		DSN string
	}
}

func newApp(conf *Config) *App {

	app := &App{
		Config: conf,
	}

	router := newRouter(app)
	server := http.Server{
		Handler:  router,
		Addr:     fmt.Sprintf(":%d", conf.Server.Port),
		ErrorLog: log.New(app.Logger, "", log.LstdFlags),
	}

	app.Server = &server

	return app
}

func (app *App) ListenAndServe() error {
	app.infos("server is starting", map[string]interface{}{"addr": app.Config.Server.Port})
	return app.Server.ListenAndServe()
}
