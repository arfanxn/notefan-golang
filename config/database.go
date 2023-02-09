package config

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func InitializeDB() (*sql.DB, error) {
	dbConnName := os.Getenv("DB_CONNECTION")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

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
