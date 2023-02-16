package routes

import (
	"database/sql"
	"net/http"

	"github.com/notefan-golang/controllers"
	"github.com/notefan-golang/repositories"

	"github.com/gorilla/mux"
)

func registerPageRoutes(subRouter *mux.Router, db *sql.DB) {
	pageRepository := repositories.NewPageRepository(db)
	pageController := controllers.NewPageController(pageRepository)

	// Page sub routes
	pages := subRouter.PathPrefix("/pages").Subrouter()
	pages.HandleFunc("", pageController.Get).Methods(http.MethodGet)
}
