package main

import (
	"barlio/cmd/server/model"
	"flag"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
)

func main() {
	conf := Config{}
	flag.IntVar(&conf.Server.Port, "addr", 80, "the server port")
	flag.IntVar(&conf.Session.Duration, "lifetime", 80, "session lifetime")
	flag.StringVar(&conf.DB.DSN, "dsn", "postgres://barlio:barliopass@localhost:5432/barlio?sslmode=disable", "database dsn")
	flag.Parse()

	app := newApp(&conf)

	log := newLogger(os.Stdout)
	app.setLog(log)

	db, err := newDB(conf.DB.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	app.setDB(db)

	app.models = &models{
		user:    &model.UserModel{DB: db},
		article: &model.ArticleModel{DB: db},
		token:   &model.TokenModel{DB: db},
	}

	templates, err := appPage()
	if err != nil {
		log.Fatal(err)
	}
	app.Templates = templates

	sessionManager := scs.New()
	sessionManager.Lifetime = time.Duration(conf.Session.Duration)
	store := postgresstore.NewWithCleanupInterval(db, time.Hour*24)
	sessionManager.Store = store
	app.SessionManager = sessionManager

	if err := app.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
