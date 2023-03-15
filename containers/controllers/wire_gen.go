// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package controllers

import (
	"database/sql"
	"github.com/notefan-golang/controllers"
	"github.com/notefan-golang/policies"
	"github.com/notefan-golang/repositories"
	"github.com/notefan-golang/services"
)

// Injectors from injector.go:

func InitializeAuthController(db *sql.DB) *controllers.AuthController {
	userRepository := repositories.NewUserRepository(db)
	tokenRepository := repositories.NewTokenRepository(db)
	mediaRepository := repositories.NewMediaRepository(db)
	authService := services.NewAuthService(userRepository, tokenRepository, mediaRepository)
	authController := controllers.NewAuthController(authService)
	return authController
}

func InitializeUserController(db *sql.DB) *controllers.UserController {
	userRepository := repositories.NewUserRepository(db)
	mediaRepository := repositories.NewMediaRepository(db)
	userService := services.NewUserService(userRepository, mediaRepository)
	userController := controllers.NewUserController(userService)
	return userController
}

func InitializeMediaController(db *sql.DB) *controllers.MediaController {
	mediaRepository := repositories.NewMediaRepository(db)
	mediaService := services.NewMediaService(mediaRepository)
	mediaController := controllers.NewMediaController(mediaService)
	return mediaController
}

func InitializeSpaceController(db *sql.DB) *controllers.SpaceController {
	spaceRepository := repositories.NewSpaceRepository(db)
	userRoleSpaceRepository := repositories.NewUserRoleSpaceRepository(db)
	roleRepository := repositories.NewRoleRepository(db)
	mediaRepository := repositories.NewMediaRepository(db)
	spaceService := services.NewSpaceService(spaceRepository, userRoleSpaceRepository, roleRepository, mediaRepository)
	permissionRepository := repositories.NewPermissionRepository(db)
	spacePolicy := policies.NewSpacePolicy(permissionRepository, userRoleSpaceRepository)
	spaceController := controllers.NewSpaceController(spaceService, spacePolicy)
	return spaceController
}
