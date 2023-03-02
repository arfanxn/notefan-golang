package singletons

import (
	"sync"

	"github.com/notefan-golang/config"
	configContainer "github.com/notefan-golang/containers/config"
)

var appSingleton *config.App

var mutex = new(sync.Mutex)

// Get the app singleton instance
func GetApp() (*config.App, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if appSingleton == nil {
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
