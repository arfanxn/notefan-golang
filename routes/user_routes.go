package routes

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/notefan-golang/containers"
)

func registerUserRoutes(router *mux.Router, db *sql.DB) {
	userController := containers.InitializeUserController(db)

	// User subrouter
	users := router.PathPrefix("/users").Subrouter()

	users.HandleFunc("/self", userController.Self).Methods(http.MethodGet)
}
