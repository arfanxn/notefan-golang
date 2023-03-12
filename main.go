package main

import (
	"net/http"

	"github.com/notefan-golang/cmd"
	"github.com/notefan-golang/containers/singletons"
	"github.com/notefan-golang/helpers/errorh"
)

func main() {
	// This will guess the env based on the command line arguments
	cmd.GuessENV()

	// Initialize the Application
	app, err := singletons.GetApp()
	errorh.LogFatal(err)

	// These functions will run when some commands are executed
	cmd.RunSeeder(app.DB)

	// Start the application server
	err = http.ListenAndServe(":8080", app.Router)
	errorh.LogFatal(err)
}
