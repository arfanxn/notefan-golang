package routes

import (
	"net/http"

	"github.com/notefan-golang/config"
	"github.com/notefan-golang/controllers"
	"github.com/notefan-golang/repositories"

	"github.com/gorilla/mux"
)

func initializePageRoutes(app *config.App, subRouter *mux.Router) {
	pageRepository := repositories.NewPageRepository(app.DB)
	pageController := controllers.NewPageController(pageRepository)

	// Page sub routes
	pages := subRouter.PathPrefix("/pages").Subrouter()
	pages.HandleFunc("", pageController.Get).Methods(http.MethodGet)
}
