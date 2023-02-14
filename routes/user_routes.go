package routes

import (
	"notefan-golang/config"
	"notefan-golang/controllers"
	"notefan-golang/repositories"
	"notefan-golang/services"

	"github.com/gorilla/mux"
)

func initializeUserRoutes(app *config.App, subRouter *mux.Router) {
	userRepository := repositories.NewUserRepository(app.DB)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)
	_ = userController

	// User sub routes
	// TODO
	// users := subRouter.PathPrefix("/users").Subrouter()
	// users.HandleFunc("", userController.Something).Methods(http.MethodGet)
}
