//go:build wireinject
// +build wireinject

package controllers

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/notefan-golang/controllers"
	"github.com/notefan-golang/repositories"
	"github.com/notefan-golang/services"
)

func InitializeAuthController(db *sql.DB) *controllers.AuthController {
	wire.Build(repositories.NewUserRepository, repositories.NewMediaRepository, services.NewAuthService, controllers.NewAuthController)
	return nil
}

func InitializeUserController(db *sql.DB) *controllers.UserController {
	wire.Build(repositories.NewMediaRepository, repositories.NewUserRepository, services.NewUserService, controllers.NewUserController)
	return nil
}
