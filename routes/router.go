package routes

import (
	"net/http"
	"notefan-golang/config"
	"notefan-golang/helper"
	"notefan-golang/middlewares"
)

func InitializeRouter(app *config.App) {
	initializeApiRouter(app)
	initializeFileServer(app)

	err := http.ListenAndServe(":8080", app.Router)
	helper.ErrorLogFatal(err)
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

func initializeFileServer(app *config.App) {
	fs := http.FileServer(http.Dir("./public"))
	app.Router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))
}
