package routes

import (
	"net/http"
	"notefan-golang/config"
	"notefan-golang/helper"
	"notefan-golang/middlewares"
)

func InitializeRouter(app *config.App) {
	initializeApiRouter(app)

	err := http.ListenAndServe(":8080", app.Router)
	helper.LogFatalIfError(err)
}

func initializeApiRouter(app *config.App) {
	/* API Subroutes */
	guestApi := app.Router.PathPrefix("/api").Subrouter()
	api := app.Router.PathPrefix("/api").Subrouter()
	api.Use(middlewares.AuthenticateMiddleware)

	/* Authentication Router */
	initializeAuthRouter(app, guestApi)

	/* Page Router */
	initializePageRouter(app, api)
}
