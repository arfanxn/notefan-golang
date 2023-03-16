package routes

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	cc "github.com/notefan-golang/containers/controllers"
)

func registerPageRoutes(subRouter *mux.Router, db *sql.DB) {
	pageController := cc.InitializePageController(db)

	// Page sub routes
	spacesIdPages := subRouter.PathPrefix("/spaces/{space_id}/pages").Subrouter()
	spacesIdPages.HandleFunc("", pageController.Get).Methods(http.MethodGet)
	spacesIdPages.HandleFunc("/{page_id}", pageController.Find).Methods(http.MethodGet)
	spacesIdPages.HandleFunc("", pageController.Create).Methods(http.MethodPost)
	spacesIdPages.HandleFunc("/{page_id}", pageController.Update).Methods(http.MethodPut)
	spacesIdPages.HandleFunc("/{page_id}", pageController.Delete).Methods(http.MethodDelete)
}
