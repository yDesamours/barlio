package main

import (
	"barlio/cmd/server/config"
	"barlio/cmd/server/router"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", "80", "the server port")
	flag.Parse()

	server := http.Server{Addr: fmt.Sprintf(":%s", *addr)}
	app := config.NewApp(&server)
	router.Router(app)

	if err := app.Server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
