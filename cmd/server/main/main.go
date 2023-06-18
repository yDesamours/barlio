package main

import (
	"barlio/cmd/server/model"
	"barlio/internal/mailer"
	"flag"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
)

func main() {
	conf := Config{}
	flag.IntVar(&conf.Server.Port, "addr", 80, "the server port")
	flag.IntVar(&conf.Session.Duration, "lifetime", 24, "session lifetime")
	flag.StringVar(&conf.DB.DSN, "dsn", "postgres://barlio:barliopass@localhost:5432/barlio?sslmode=disable", "database dsn")
	flag.StringVar(&conf.SMTP.Host, "smtphost", "sandbox.smtp.mailtrap.io", "smtp server host")
	flag.StringVar(&conf.SMTP.Username, "smtpuser", "73a9a7e70105ed", "smtp server username")
	flag.StringVar(&conf.SMTP.Password, "smtppass", "bf3b98d194e0d8", "smtp server password")
	flag.IntVar(&conf.SMTP.Port, "smtpport", 25, "smtp server port")
	flag.StringVar(&conf.SMTP.Sender, "smtp sender", "yveltamours@gmail.com", "mail sender")
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

	pageTemplates, err := appPage()
	if err != nil {
		log.Fatal(err)
	}
	app.PageTemplates = pageTemplates

	mailTemplates, err := mailTemplates()
	if err != nil {
		log.Fatal(err)
	}
	app.MailTemplate = mailTemplates

	sessionManager := scs.New()
	sessionManager.Lifetime = time.Duration(conf.Session.Duration) * time.Hour
	sessionManager.Store = postgresstore.NewWithCleanupInterval(db, time.Hour*24)
	app.SessionManager = sessionManager

	smtp := &conf.SMTP
	mailer := mailer.NewMailer(smtp.Port, smtp.Host, smtp.Username, smtp.Password, smtp.Sender)
	app.Mailer = mailer

	if err := app.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
