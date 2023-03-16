package routes

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	cc "github.com/notefan-golang/containers/controllers"
)

// registerSpaceMemberRoutes registers routes for Space member module
func registerSpaceMemberRoutes(router *mux.Router, db *sql.DB) {
	spaceMemberController := cc.InitializeSpaceMemberController(db)

	// Space subrouters
	spacesIdMembers := router.PathPrefix("/spaces/{space_id}/members").Subrouter()

	spacesIdMembers.HandleFunc("", spaceMemberController.Get).Methods(http.MethodGet)
	spacesIdMembers.HandleFunc("/{member_id}", spaceMemberController.Find).Methods(http.MethodGet)
	spacesIdMembers.HandleFunc("", spaceMemberController.Invite).Methods(http.MethodPost)
	spacesIdMembers.HandleFunc("/{member_id}/role", spaceMemberController.UpdateRole).Methods(http.MethodPut)
	spacesIdMembers.HandleFunc("/{member_id}", spaceMemberController.Remove).Methods(http.MethodDelete)
}
