package main

import (
	"barlio/cmd/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

type App struct {
	Server         *http.Server
	SessionManager *scs.SessionManager
	DB             *sql.DB
	Logger         *Logger
	Config         Config
	User           *models.UserModel
	Article        *models.ArticleModel
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

	server := http.Server{
		Addr: fmt.Sprintf(":%d", conf.Server.Port),
	}

	db, err := newDB(conf.DB.DSN)
	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		DB:             db,
		Config:         *conf,
		Server:         &server,
		SessionManager: sessionManager,
		Logger:         newLogger(os.Stdout),
		User:           &models.UserModel{DB: db},
		Article:        &models.ArticleModel{DB: db},
	}
	app.routes()

	return app
}

func newDB(dsn string) (*sql.DB, error) {
	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = DB.Ping()
	if err != nil {
		return nil, err
	}

	return DB, nil
}

func (app *App) ListenAndServe() error {
	app.Logger.WriteInfos("server is starting", map[string]interface{}{"addr": app.Config.Server.Port})
	return app.Server.ListenAndServe()
}
