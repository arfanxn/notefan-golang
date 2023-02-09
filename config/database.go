package config

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func InitializeDB() (*sql.DB, error) {
	dbConnName, err := GetENVKey("DB_CONNECTION")
	dbHost, err := GetENVKey("DB_HOST")
	dbPort, err := GetENVKey("DB_PORT")
	dbName, err := GetENVKey("DB_DATABASE")
	dbUsername, err := GetENVKey("DB_USERNAME")
	dbPassword, err := GetENVKey("DB_PASSWORD")
	if err != nil {
		return &sql.DB{}, err
	}

	db, err := sql.Open(dbConnName, dbUsername+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbName)
	if err != nil {
		return &sql.DB{}, err
	}

	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	/*
	* Run this for database migration
	* migrate -database "mysql://root@tcp(localhost:3306)/notefan" -path database/migrations up
	 */

	return db, err
}
