package routes

import (
	"net/http"
	"notefan-golang/config"
	"notefan-golang/controllers"
	"notefan-golang/repositories"

	"github.com/gorilla/mux"
)

func initializePageRoutes(app *config.App, subRouter *mux.Router) {
	pageRepository := repositories.NewPageRepository(app.DB)
	pageController := controllers.NewPageController(pageRepository)

	// Page sub routes
	pages := subRouter.PathPrefix("/pages").Subrouter()
	pages.HandleFunc("", pageController.Get).Methods(http.MethodGet)
}
