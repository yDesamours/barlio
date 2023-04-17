package config

import (
	"net/http"
	"os"
	"path/filepath"
)

type App struct {
	Server *http.Server
	Path   string
}

func NewApp(server *http.Server) *App {
	exe, err := os.Executable()
	if err != nil {
		return nil
	}

	path := filepath.Dir(exe)
	return &App{Server: server, Path: path}
}
