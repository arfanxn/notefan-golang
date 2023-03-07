package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/notefan-golang/controllers"
)

func registerHandlerRoutes(router *mux.Router, db *sql.DB) {
	router.NotFoundHandler = router.NewRoute().
		HandlerFunc(controllers.NewNotFoundController().HandlerFunc).
		GetHandler()
}
