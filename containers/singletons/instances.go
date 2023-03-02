package singletons

import (
	"sync"

	"github.com/notefan-golang/config"
	configContainer "github.com/notefan-golang/containers/config"
	"github.com/notefan-golang/helpers/cmdh"
)

var appSingleton *config.App

var mutex = new(sync.Mutex)

// Get the app singleton instance
func GetApp() (*config.App, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if appSingleton == nil {
		switch true {
		case cmdh.UserFirstArgIs("test"):
			config.LoadTestENV()
			break
		default:
			config.LoadENV()
		}

		app, err := configContainer.InitializeApp()
		if err != nil {
			return nil, err
		}

		appSingleton = app
		return appSingleton, nil

	} else {
		return appSingleton, nil
	}
}
