package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/notefan-golang/middlewares"
)

// Register main middlewares for all routes
func registerMainMiddlewares(router *mux.Router) {
	router.Use(
		middlewares.RecoveryMiddleware,
		middlewares.FormDataMiddleware,
		middlewares.NeuterMiddleware,
	)
}

// InitializeRouter
func InitializeRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	registerMainMiddlewares(router)

	registerApiRoutes(router, db)
	registerFileServer(router)

	return router
}
