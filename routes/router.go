package routes

import (
	"net/http"
	"notefan-golang/config"
	"notefan-golang/helper"
	"notefan-golang/middlewares"
)

// Register main middlewares for all routes
func registerMainMiddlewares(app *config.App) {
	app.Router.Use(middlewares.RecoveryMiddleware)
}

func InitializeRouter(app *config.App) {
	registerMainMiddlewares(app)

	initializeApiRoutes(app)
	initializeFileServer(app)

	err := http.ListenAndServe(":8080", app.Router)
	helper.ErrorLogFatal(err)
}

func initializeApiRoutes(app *config.App) {
	// Prefix
	apiPathPrefix := "/api"

	// API Subroutes
	guestApi := app.Router.PathPrefix(apiPathPrefix).Subrouter()
	api := app.Router.PathPrefix(apiPathPrefix).Subrouter()
	api.Use(middlewares.AuthenticateMiddleware)

	// Authentication Routes
	initializeAuthRoutes(app, guestApi)

	// Page Routes
	initializePageRoutes(app, api)
}

func initializeFileServer(app *config.App) {
	fileServer := http.FileServer(http.Dir("./public")) // make file server and set the root directory
	app.Router.PathPrefix("/public/").
		Handler(http.StripPrefix("/public/", fileServer)) // register file server to router
}
