package routes

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/notefan-golang/containers"
	"github.com/notefan-golang/middlewares"
)

func registerAuthRoutes(router *mux.Router, db *sql.DB) {
	authController := containers.InitializeAuthController(db)

	// Auth subrouters
	users := router.PathPrefix("/users").Subrouter()
	usersLogout := users.PathPrefix("/logout").Subrouter()
	usersLogout.Use(middlewares.AuthenticateMiddleware)

	// Login and register routes
	users.HandleFunc("/login", authController.Login).Methods(http.MethodPost)
	users.HandleFunc("/register", authController.Register).Methods(http.MethodPost)

	// Logout Routes
	usersLogout.HandleFunc("", authController.Logout).Methods(http.MethodDelete)
}
