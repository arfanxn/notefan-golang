package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/notefan-golang/config"
	"github.com/notefan-golang/database/seeders"
	"github.com/notefan-golang/helper"

	"github.com/golang-migrate/migrate/v4"
)

func main() {
	// Guess environment
	guessENV()

	// Initialize the Application
	app, err := InitializeApp()
	helper.ErrorLogPanic(err)

	// These functions will run when some commands are executed
	runSeeder(app.DB)
	runPlayground()

	// Start the application server
	err = http.ListenAndServe(":8080", app.Router)
	helper.ErrorLogPanic(err)
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

// guessENV will guess the environment variable is it on production or development or test
func guessENV() {
	switch true {
	case helper.CMDUserFirstArgIs("test"):
		config.LoadTestENV()
		break
	default:
		config.LoadENV()
	}
}
