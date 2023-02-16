package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

// registerFileServer registers file server to the router
func registerFileServer(router *mux.Router) {
	fileServer := http.
		FileServer(http.Dir("./public")) // make file server and set the root directory
	router.PathPrefix("/public/").
		Handler(http.StripPrefix("/public/", fileServer)) // register file server to router
}
