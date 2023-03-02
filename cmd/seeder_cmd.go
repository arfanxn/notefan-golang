package cmd

import (
	"database/sql"

	"github.com/notefan-golang/database/seeders"
	"github.com/notefan-golang/helpers/cmdh"
)

// RunSeeder check if the command first argument equals to "seed"
// if its, it will run database seeder
func RunSeeder(db *sql.DB) {
	if !cmdh.UserFirstArgIs("seed") {
		return
	}

	seeder := seeders.NewDatabaseSeeder(db)
	seeder.Run()
}
