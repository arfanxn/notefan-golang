package routes

import (
	"net/http"
	"notion-golang/config"
	"notion-golang/controllers"
	"notion-golang/repositories"
	"notion-golang/services"

	"github.com/gorilla/mux"
)

func initializeAuthRouter(app *config.App, subRouter *mux.Router) {
	userRepo := repositories.NewUserRepo(app.DBTX)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	subRouter.HandleFunc("/login", authController.Login).Methods(http.MethodGet)
	subRouter.HandleFunc("/logout", authController.Logout).Methods(http.MethodDelete)
	subRouter.HandleFunc("/register", authController.Register).Methods(http.MethodPost)
}
