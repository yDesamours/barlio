package main

import (
	"flag"
	"log"
)

func main() {
	var conf Config
	flag.IntVar(&conf.Server.Port, "addr", 80, "the server port")
	flag.Parse()

	app := newApp(&conf)

	if err := app.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
