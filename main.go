package main

import (
	"database/sql"
	"notefan-golang/config"
	"notefan-golang/database/seeders"
	"notefan-golang/helper"
	"notefan-golang/routes"
	"os"

	"github.com/golang-migrate/migrate/v4"
)

func main() {
	// Initialize the Application
	app := config.InitializeApp()

	runSeeder(app.DB)
	runPlayground()

	// Initialize routes of the application
	routes.InitializeRoutes(app)
}

// runPlayground check if the command first argument equals to "seed"
// if its, it will run database seeder
func runSeeder(db *sql.DB) {
	if !helper.CMDUserFirstArgIs("seed") {
		return
	}

	seeder := seeders.NewDatabaseSeeder(db)
	seeder.Run()
}

// runPlayground check if the command first argument equals to "play"
// if its, it will run program as a playground
func runPlayground() {
	if !helper.CMDUserFirstArgIs("play") {
		return
	}
	defer os.Exit(0)

	m, err := migrate.New("database/migrations", "mysql://root@tcp(localhost:3306)/notefan_test")
	helper.ErrorPanic(err)
	m.Run()
}
