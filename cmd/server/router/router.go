package router

import (
	"barlio/cmd/server/config"
	"barlio/cmd/server/handler"
	"net/http"
)

func routes(app *config.App) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/home", handler.Home(app))
	router.HandleFunc("/", handler.NotFund(app))

	return router
}
func Router(app *config.App) error {
	router := routes(app)
	app.Server.Handler = router

	return nil
}
