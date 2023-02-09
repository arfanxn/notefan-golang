package routes

import (
	"net/http"
	"notefan-golang/config"
	"notefan-golang/controllers"
	"notefan-golang/repositories"

	"github.com/gorilla/mux"
)

func initializePageRouter(app *config.App, subRouter *mux.Router) {
	pageRepo := repositories.NewPageRepo(app.DBTX)
	pageController := controllers.NewPageController(pageRepo)

	// Page sub routes
	pages := subRouter.PathPrefix("/pages").Subrouter()
	pages.HandleFunc("/", pageController.Get).Methods(http.MethodGet)
}
