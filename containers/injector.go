//go:build wireinject
// +build wireinject

package containers

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/notefan-golang/config"
	"github.com/notefan-golang/controllers"
	"github.com/notefan-golang/repositories"
	"github.com/notefan-golang/services"

	"github.com/google/wire"
)

func InitializeApp() (*config.App, error) {
	wire.Build(config.InitializeDB, mux.NewRouter, config.NewApp)
	return nil, nil
}

func InitializeAuthController(db *sql.DB) *controllers.AuthController {
	wire.Build(repositories.NewUserRepository, services.NewAuthService, controllers.NewAuthController)
	return nil
}

func InitializeUserController(db *sql.DB) *controllers.UserController {
	wire.Build(repositories.NewUserRepository, services.NewUserService, controllers.NewUserController)
	return nil
}
