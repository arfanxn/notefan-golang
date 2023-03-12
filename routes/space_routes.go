package routes

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	cc "github.com/notefan-golang/containers/controllers"
)

// registerSpaceRoutes registers routes for user module
func registerSpaceRoutes(router *mux.Router, db *sql.DB) {
	spaceController := cc.InitializeSpaceController(db)

	// Space subrouters
	usersSelf := router.PathPrefix("/users/self").Subrouter()
	usersSelfSpaces := usersSelf.PathPrefix("/spaces").Subrouter()
	spaces := router.PathPrefix("/spaces").Subrouter()

	usersSelfSpaces.HandleFunc("", spaceController.Get).Methods(http.MethodGet)
	spaces.HandleFunc("/{id}", spaceController.Find).Methods(http.MethodGet)
	spaces.HandleFunc("", spaceController.Create).Methods(http.MethodPost)
	spaces.HandleFunc("/{id}", spaceController.Update).Methods(http.MethodPut)
	spaces.HandleFunc("/{id}", spaceController.Delete).Methods(http.MethodDelete)
}
