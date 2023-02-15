package integrations

import (
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/notefan-golang/helper"
)

func TestMain(m *testing.M) {
	setup()
	defer teardown()
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
