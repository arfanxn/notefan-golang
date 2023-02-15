package routes

import (
	"net/http"
	"notefan-golang/config"
	"notefan-golang/containers"
	"notefan-golang/middlewares"

	"github.com/gorilla/mux"
)

func initializeAuthRoutes(app *config.App, subRouter *mux.Router) {
	authController := containers.InitializeAuthController(app.DB)

	// Login and register routes
	users := subRouter.PathPrefix("/users").Subrouter()
	users.HandleFunc("/login", authController.Login).Methods(http.MethodPost)
	users.HandleFunc("/register", authController.Register).Methods(http.MethodPost)

	// Logout Route
	usersLogout := users.PathPrefix("/logout").Subrouter()
	usersLogout.Use(middlewares.AuthenticateMiddleware)
	usersLogout.HandleFunc("", authController.Logout).Methods(http.MethodDelete)
}
