package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/notefan-golang/config"
	"github.com/notefan-golang/middlewares"
)

// Register main middlewares for all routes
func registerMainMiddlewares(app *config.App) {
	app.Router.Use(middlewares.RecoveryMiddleware)
}

// InitializeRoutes
func InitializRoutes(app *config.App) *mux.Router {
	registerMainMiddlewares(app)

	initializeApiRoutes(app)
	initializeFileServer(app)

	return app.Router
}

// initializeFileServer initializes the file server router
func initializeFileServer(app *config.App) {
	fileServer := http.FileServer(http.Dir("./public")) // make file server and set the root directory
	app.Router.PathPrefix("/public/").
		Handler(http.StripPrefix("/public/", fileServer)) // register file server to router
}
