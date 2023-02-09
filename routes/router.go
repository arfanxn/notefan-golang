package routes

import (
	"net/http"
	"notefan-golang/config"
	"notefan-golang/helper"
)

func InitializeRouter(app *config.App) {
	initializeApiRouter(app)

	err := http.ListenAndServe(":8080", app.Router)
	helper.LogFatalIfError(err)
}

func initializeApiRouter(app *config.App) {
	/* API Subrouter */
	api := app.Router.PathPrefix("/api").Subrouter()

	/* Page Router */
	initializePageRouter(app, api)

	/* Authentication Router */
	initializeAuthRouter(app, api)
}
