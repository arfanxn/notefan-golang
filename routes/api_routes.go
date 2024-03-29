package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/notefan-golang/middlewares"
)

func registerApiRoutes(router *mux.Router, db *sql.DB) {
	// Prefix
	apiPathPrefix := "/api"

	// API Subrouters
	publicApi := router.PathPrefix(apiPathPrefix).Subrouter()
	protectedApi := router.PathPrefix(apiPathPrefix).Subrouter()
	protectedApi.Use(middlewares.AuthenticateMiddleware)

	// Authentication Routes
	registerAuthRoutes(publicApi, db)

	// User Routes
	registerUserRoutes(protectedApi, db)

	// Media Routes
	registerMediaRoutes(protectedApi, db)

	// Space Routes
	registerSpaceRoutes(protectedApi, db)

	// Space Member Routes
	registerSpaceMemberRoutes(protectedApi, db)

	// Page Routes
	registerPageRoutes(protectedApi, db)

	// Page Content Routes
	registerPageContentRoutes(protectedApi, db)
}
