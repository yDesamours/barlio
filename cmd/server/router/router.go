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
	router.Handle("/statics/", http.StripPrefix("/statics/", handler.FileServer(app)))

	return router
}

func Router(app *config.App) error {
	router := routes(app)
	app.Server.Handler = recoverMiddleware(router)

	return nil
}

func recoverMiddleware(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				http.Error(w, "an error has occured", http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	}
}
