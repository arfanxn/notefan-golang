package integrations

import (
	"net/http"
	"os"

	"net/http/cookiejar"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/notefan-golang/helper"
)

// httpClient is a client for interacting with the integration testing
var httpClient *http.Client = func() *http.Client {
	cookiejar, err := cookiejar.New(nil)
	helper.ErrorLogPanic(err)
	return &http.Client{
		Timeout: time.Second * 3,
		Jar:     cookiejar,
	}
}()

func TestMain(m *testing.M) {
	setup()
	defer func() {
		teardown()
	}()

	// Run tests
	m.Run()
}

// setup sets up the test
func setup() {
	// this will check if migration up fails perhaps it coz of "no changes" error so it should drop and up again to achieve migration up successfully
	if migrateDB().Up() != nil {
		helper.ErrorLogPanic(migrateDB().Drop())
		helper.ErrorLogPanic(migrateDB().Up())
	}
}

// teardown teardowns the test
func teardown() {
	time.Sleep(time.Second * 10)
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
