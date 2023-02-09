package config

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func InitializeDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/notion")

	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	/*
	* Run this for database migration
	* migrate -database "mysql://root@tcp(localhost:3306)/notion" -path database/migrations up
	 */

	return db, err
}
