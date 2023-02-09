package config

import (
	"database/sql"

	"github.com/gorilla/mux"
)

type App struct {
	DBTX   *sql.DB
	Router *mux.Router
}

func NewApp(dbtx *sql.DB, router *mux.Router) *App {
	return &App{
		DBTX:   dbtx,
		Router: router,
	}
}
