package main

import (
	"fmt"
	"net/http"
)

func (app App) recoverMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println(err)
				http.Error(w, "an error has occured", http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
