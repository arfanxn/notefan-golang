package tests

import (
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/notefan-golang/config"
	"github.com/notefan-golang/helpers/errorh"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// httpClient is a client like web browser that can be used to communicate to the server
var httpClient *http.Client = nil

// GetHTTPClient returns configured http.Client for testing purposes
func GetHTTPClient() *http.Client {
	if httpClient != nil {
		return httpClient
	}

	cookiejar, err := cookiejar.New(nil)
	errorh.LogFatal(err)

	httpClient = &http.Client{
		Timeout: time.Second * 3,
		Jar:     cookiejar,
	}
	return httpClient
}

// Setup sets up the test
func Setup() {
	config.LoadTestENV()

	// this will check if migration up fails perhaps it coz of "no changes" error so it should drop and up again to achieve migration up successfully
	if migrateDB().Up() != nil {
		errorh.LogFatal(migrateDB().Drop())
		errorh.LogFatal(migrateDB().Up())
	}
}

// Teardown teardowns the test
func Teardown() {
	errorh.LogFatal(migrateDB().Drop())
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
	errorh.LogFatal(err)
	return m
}
