package routes

import (
	"net/http"
	"notefan-golang/config"
	"notefan-golang/controllers"
	"notefan-golang/middlewares"
	"notefan-golang/repositories"
	"notefan-golang/services"

	"github.com/gorilla/mux"
)

func initializeAuthRoutes(app *config.App, subRouter *mux.Router) {
	// Preapare the dependencies
	userRepo := repositories.NewUserRepo(app.DB)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	// Login and register routes
	users := subRouter.PathPrefix("/users").Subrouter()
	users.HandleFunc("/login", authController.Login).Methods(http.MethodPost)
	users.HandleFunc("/register", authController.Register).Methods(http.MethodPost)

	// Logout Route
	usersLogout := users.PathPrefix("/logout").Subrouter()
	usersLogout.Use(middlewares.AuthenticateMiddleware)
	usersLogout.HandleFunc("", authController.Logout).Methods(http.MethodDelete)
}
