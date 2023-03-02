//go:build wireinject
// +build wireinject

package config

import (
	"github.com/notefan-golang/config"
	"github.com/notefan-golang/routes"

	"github.com/google/wire"
)

func InitializeApp() (*config.App, error) {
	wire.Build(config.InitializeDB, routes.InitializeRouter, config.NewApp)
	return nil, nil
}
