package tests

import (
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/notefan-golang/config"
	"github.com/notefan-golang/helper"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// app is wrapper for application configuration
var app *config.App = nil

// GetApp returns the application configuration
func GetApp() *config.App {
	if app != nil {
		return app
	}

	a, err := InitializeApp()
	if err != nil {
		helper.ErrorLogPanic(err)
		return nil
	}
	return a
}

// httpClient is a client like web browser that can be used to communicate to the server
var httpClient *http.Client = nil

// GetHTTPClient returns configured http.Client for testing purposes
func GetHTTPClient() *http.Client {
	if httpClient != nil {
		return httpClient
	}

	cookiejar, err := cookiejar.New(nil)
	helper.ErrorLogPanic(err)

	httpClient = &http.Client{
		Timeout: time.Second * 3,
		Jar:     cookiejar,
	}
	return httpClient
}

// Setup sets up the test
func Setup() {
	// this will check if migration up fails perhaps it coz of "no changes" error so it should drop and up again to achieve migration up successfully
	if migrateDB().Up() != nil {
		helper.ErrorLogPanic(migrateDB().Drop())
		helper.ErrorLogPanic(migrateDB().Up())
	}
}

// Teardown teardowns the test
func Teardown() {
	helper.ErrorLogPanic(migrateDB().Drop())
}

func migrateDB() *migrate.Migrate {
	dbConnName := os.Getenv("DB_CONNECTION") // mysql
	dbHost := os.Getenv("DB_HOST")           // 8080
	dbPort := os.Getenv("DB_PORT")           // localhost
	dbName := os.Getenv("DB_DATABASE")       // notefan
	dbUsername := os.Getenv("DB_USERNAME")   // root
	dbPassword := os.Getenv("DB_PASSWORD")   // password
	m, err := migrate.New(
		"file://database/migrations",
		dbConnName+"://"+dbUsername+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?parseTime=true",
	)
	helper.ErrorPanic(err)
	return m
}
