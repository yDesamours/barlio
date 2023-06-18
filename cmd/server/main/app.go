package main

import (
	"barlio/cmd/server/model"
	"barlio/internal/mailer"
	"fmt"
	"log"

	"database/sql"
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
	PageTemplates  templateMap
	Mailer         *mailer.Mailer
	MailTemplate   templateMap
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
	Session struct {
		Duration int
	}
	SMTP struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
}

func newApp(conf *Config) *App {

	app := &App{
		Config: conf,
	}

	return app
}

func (app *App) ListenAndServe() error {
	router := newRouter(app)
	server := http.Server{
		Handler:  router,
		Addr:     fmt.Sprintf(":%d", app.Config.Server.Port),
		ErrorLog: log.New(app.Logger, "", log.LstdFlags),
	}

	app.Server = &server
	app.infos("server is starting", map[string]interface{}{"addr": app.Config.Server.Port})
	return app.Server.ListenAndServe()
}
