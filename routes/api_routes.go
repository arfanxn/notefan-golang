package routes

import (
	"github.com/notefan-golang/config"
	"github.com/notefan-golang/middlewares"
)

func initializeApiRoutes(app *config.App) {
	// Prefix
	apiPathPrefix := "/api"

	// API Subroutes
	publicApi := app.Router.PathPrefix(apiPathPrefix).Subrouter()
	protectedApi := app.Router.PathPrefix(apiPathPrefix).Subrouter()
	protectedApi.Use(middlewares.AuthenticateMiddleware)

	// Authentication Routes
	initializeAuthRoutes(app, publicApi)

	// User Routes
	initializeUserRoutes(app, protectedApi)

	// Page Routes
	initializePageRoutes(app, protectedApi)
}
