// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package containers

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/notefan-golang/config"
	"github.com/notefan-golang/controllers"
	"github.com/notefan-golang/repositories"
	"github.com/notefan-golang/services"
)

// Injectors from injector.go:

func InitializeApp() (*config.App, error) {
	db, err := config.InitializeDB()
	if err != nil {
		return nil, err
	}
	router := mux.NewRouter()
	app := config.NewApp(db, router)
	return app, nil
}

func InitializeAuthController(db *sql.DB) *controllers.AuthController {
	userRepository := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepository)
	authController := controllers.NewAuthController(authService)
	return authController
}
