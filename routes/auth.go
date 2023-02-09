package routes

import (
	"net/http"
	"notefan-golang/config"
	"notefan-golang/controllers"
	"notefan-golang/repositories"
	"notefan-golang/services"

	"github.com/gorilla/mux"
)

func initializeAuthRouter(app *config.App, subRouter *mux.Router) {
	userRepo := repositories.NewUserRepo(app.DBTX)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	// User sub routes
	users := subRouter.PathPrefix("/users").Subrouter()
	users.HandleFunc("/login", authController.Login).Methods(http.MethodPost)
	users.HandleFunc("/logout", authController.Logout).Methods(http.MethodDelete)
	users.HandleFunc("/register", authController.Register).Methods(http.MethodPost)
}
