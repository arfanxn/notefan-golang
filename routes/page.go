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

	subRouter.HandleFunc("/pages", pageController.Get).Methods(http.MethodGet)
}
