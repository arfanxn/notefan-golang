package routes

import (
	"net/http"
	"notion-golang/config"
	"notion-golang/controllers"
	"notion-golang/repositories"

	"github.com/gorilla/mux"
)

func initializePageRouter(app *config.App, subRouter *mux.Router) {
	pageRepo := repositories.NewPageRepo(app.DBTX)
	pageController := controllers.NewPageController(pageRepo)

	subRouter.HandleFunc("/pages", pageController.Get).Methods(http.MethodGet)
}
