//go:build wireinject
// +build wireinject

package containers

import (
	"database/sql"
	"notefan-golang/controllers"
	"notefan-golang/repositories"
	"notefan-golang/services"

	"github.com/google/wire"
)

func InitializeAuthController(db *sql.DB) *controllers.AuthController {
	wire.Build(repositories.NewUserRepository, services.NewAuthService, controllers.NewAuthController)
	return nil
}
