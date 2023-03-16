package config

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func InitializeDB() (*sql.DB, error) {
	dbConnName := os.Getenv("DB_CONNECTION") // mysql
	dbHost := os.Getenv("DB_HOST")           // 8080
	dbPort := os.Getenv("DB_PORT")           // localhost
	dbName := os.Getenv("DB_DATABASE")       // notefan
	dbUsername := os.Getenv("DB_USERNAME")   // root
	dbPassword := os.Getenv("DB_PASSWORD")   // password

	db, err := sql.Open(
		dbConnName, dbUsername+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?parseTime=true")
	if err != nil {
		return &sql.DB{}, err
	}

	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	/*
	* Run this for database migration
	* clear ; migrate -database "mysql://root@tcp(localhost:3306)/notefan" -path database/migrations drop ;
	 migrate -database "mysql://root@tcp(localhost:3306)/notefan" -path database/migrations up ; go run . seed
	 */

	return db, err
}
