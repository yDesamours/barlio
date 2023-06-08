package main

import (
	"flag"
	"os"
)

func main() {
	conf := Config{}
	flag.IntVar(&conf.Server.Port, "addr", 80, "the server port")
	flag.Parse()

	app := newApp(&conf)

	log := newLogger(os.Stdout)
	app.setLog(log)

	// db, err := newDB(conf.DB.DSN)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()
	// app.setDB(db)

	templates, err := appPage()
	if err != nil {
		log.Fatal(err)
	}
	app.Templates = templates

	if err := app.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
