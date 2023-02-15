package config

import (
	"database/sql"

	"github.com/gorilla/mux"
)

type App struct {
	DB     *sql.DB
	Router *mux.Router
}

func NewApp(DB *sql.DB, router *mux.Router) *App {
	return &App{
		DB:     DB,
		Router: router,
	}
}
