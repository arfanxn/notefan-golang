package routes

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	cc "github.com/notefan-golang/containers/controllers"
)

// registerUserRoutes registers routes for user module
func registerUserRoutes(router *mux.Router, db *sql.DB) {
	userController := cc.InitializeUserController(db)

	// User subrouters
	users := router.PathPrefix("/users").Subrouter()
	usersSelf := users.PathPrefix("/self").Subrouter()

	usersSelf.HandleFunc("", userController.Self).Methods(http.MethodGet)
	usersSelf.HandleFunc("/update", userController.Update).Methods(http.MethodPut)
}
