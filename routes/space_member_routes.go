package routes

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	cc "github.com/notefan-golang/containers/controllers"
)

// registerSpaceMemberRoutes registers routes for user module
func registerSpaceMemberRoutes(router *mux.Router, db *sql.DB) {
	spaceMemberController := cc.InitializeSpaceMemberController(db)

	// Space subrouters
	spaces := router.PathPrefix("spaces").Subrouter()

	spaces.HandleFunc("/members/{id}", spaceMemberController.Get).Methods(http.MethodGet)
	// spacesIdMembers.HandleFunc("/{member_id}", spaceMemberController.Find).Methods(http.MethodGet)
	// spaces.HandleFunc("/{id}", spaceMemberController.Update).Methods(http.MethodPut)
	// spaces.HandleFunc("/{id}", spaceMemberController.Delete).Methods(http.MethodDelete)
}
