package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

type App struct {
	Server         *http.Server
	SessionManager *scs.SessionManager
	DB             *sql.DB
	Logger         Logger
	Config         Config
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
	sessionManager := scs.New()
	sessionManager.Lifetime = time.Hour

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", conf.Server.Port),
	}

	app := &App{
		Config:         *conf,
		Server:         server,
		SessionManager: sessionManager,
		Logger:         *newLogger(os.Stdout),
	}
	app.routes()

	return app
}

func (app *App) ListenAndServe() error {
	app.Logger.WriteInfos("server is starting", map[string]interface{}{"addr": app.Config.Server.Port})
	return app.Server.ListenAndServe()
}
