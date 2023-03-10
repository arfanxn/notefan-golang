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

func InitializeMediaController(db *sql.DB) *controllers.MediaController {
	wire.Build(repositories.NewMediaRepository, services.NewMediaService, controllers.NewMediaController)
	return nil
}

func InitializeSpaceController(db *sql.DB) *controllers.SpaceController {
	wire.Build(
		repositories.NewSpaceRepository,
		repositories.NewUserRoleSpaceRepository,
		repositories.NewRoleRepository,
		repositories.NewMediaRepository,
		services.NewSpaceService,
		controllers.NewSpaceController,
	)
	return nil
}
