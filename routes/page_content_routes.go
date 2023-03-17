package routes

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	cc "github.com/notefan-golang/containers/controllers"
)

func registerPageContentRoutes(subRouter *mux.Router, db *sql.DB) {
	pageContentController := cc.InitializePageContentController(db)

	// Page Content sub routes
	pagesIdContents := subRouter.PathPrefix("/pages/{page_id}/contents").Subrouter()
	pagesIdContents.HandleFunc("", pageContentController.Get).Methods(http.MethodGet)
	pagesIdContents.HandleFunc("/{page_content_id}", pageContentController.Find).Methods(http.MethodGet)
	pagesIdContents.HandleFunc("", pageContentController.Create).Methods(http.MethodPost)
	pagesIdContents.HandleFunc("/{page_content_id}", pageContentController.Update).Methods(http.MethodPut)
	pagesIdContents.HandleFunc("/{page_content_id}", pageContentController.Delete).Methods(http.MethodDelete)
}
