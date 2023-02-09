package main

import (
	"notion-golang/config"
	"notion-golang/helper"
	"notion-golang/routes"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the environment file
	err := config.InitializeENV()
	helper.LogFatalIfError(err)

	// Initialize the Database connection
	db, err := config.InitializeDB()
	helper.LogFatalIfError(err)

	// Instantiate the router
	router := mux.NewRouter()

	// Instantiate the Application
	app := config.NewApp(db, router)

	// Initialize the Router
	routes.InitializeRouter(app)
}
