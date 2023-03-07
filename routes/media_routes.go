package routes

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	cc "github.com/notefan-golang/containers/controllers"
)

// registerMediaRoutes registers routes for user module
func registerMediaRoutes(router *mux.Router, db *sql.DB) {
	mediaController := cc.InitializeMediaController(db)

	// Media subrouters
	medias := router.PathPrefix("/medias").Subrouter()

	medias.HandleFunc("/{id}", mediaController.Find).Methods(http.MethodGet)
	medias.HandleFunc("/{id}", mediaController.Update).Methods(http.MethodPut)
	medias.HandleFunc("/{id}", mediaController.Delete).Methods(http.MethodDelete)
}
