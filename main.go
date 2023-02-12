package main

import (
	"database/sql"
	"notefan-golang/config"
	"notefan-golang/database/seeders"
	"notefan-golang/helper"
	"notefan-golang/routes"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the environment file
	err := config.InitializeENV()
	helper.LogFatalIfError(err)

	// Initialize the Database connection
	db, err := config.InitializeDB()
	helper.LogFatalIfError(err)
	seedIfNeeded(db)

	// Instantiate the router
	router := mux.NewRouter()

	// Instantiate the Application
	app := config.NewApp(db, router)

	// Initialize the Router
	routes.InitializeRouter(app)
}

// This function will check if the command arguments contains "seed"
// if its contains it will run database seeder
func seedIfNeeded(db *sql.DB) {
	if (len(os.Args) > 1) && strings.Contains(os.Args[1], "seed") {
		seeder := seeders.NewDatabaseSeeder(db)
		seeder.Run()
	}
}
