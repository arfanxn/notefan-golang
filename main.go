package main

import (
	"database/sql"
	"fmt"
	"notefan-golang/config"
	"notefan-golang/database/seeders"
	"notefan-golang/helper"
	"notefan-golang/routes"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	runPlayground()

	// Initialize the environment file
	err := config.InitializeENV()
	helper.ErrorLogFatal(err)

	// Initialize the Database connection
	db, err := config.InitializeDB()
	helper.ErrorLogFatal(err)
	runSeeder(db)

	// Instantiate the router
	router := mux.NewRouter()

	// Instantiate the Application
	app := config.NewApp(db, router)

	// Initialize the Router
	routes.InitializeRouter(app)
}

// runSeeder will check if the command first argument contains "seed"
// if its contains it will run database seeder
func runSeeder(db *sql.DB) {
	if (len(os.Args) > 1) && strings.Contains(os.Args[1], "seed") {
		seeder := seeders.NewDatabaseSeeder(db)
		seeder.Run()
	}
}

// runPlayground check if the command first argument contains "play"
// if its contains it will run program as a playground
func runPlayground() {
	if (len(os.Args) > 1) && strings.Contains(os.Args[1], "play") {
		f, _ := helper.FileRandFromDir("./controllers")
		fmt.Println(f.Stat())
	}
}
