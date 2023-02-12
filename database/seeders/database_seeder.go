package seeders

import (
	"database/sql"
	"fmt"
	"notefan-golang/helper"
	"os"
	"time"
)

type DatabaseSeeder struct {
	db      *sql.DB
	Seeders []SeederContract
}

func NewDatabaseSeeder(db *sql.DB) *DatabaseSeeder {
	return &DatabaseSeeder{db: db}
}

func (seeder *DatabaseSeeder) Run() {
	// Consoler
	fmt.Println("Running Seeder...")

	db := seeder.db

	// Inject entity seeders into strcut field
	seeder.Seeders = append([]SeederContract{
		// User and related seeders
		NewUserSeeder(db),
		NewPermissionSeeder(db),
		NewRoleSeeder(db),
		NewPermissionRoleSeeder(db),

		// Space and related seeders
		NewSpaceSeeder(db),
		NewUserRoleSpaceSeeder(db),

		// Page and related seeders
	}, seeder.Seeders...)

	for _, entitySeeder := range seeder.Seeders {
		entitySeeder.Run()
	}

	defer func() {
		fmt.Println("Seeding completed successfully")
		os.Exit(0)
	}()
}

func printStartRunning(pc uintptr) {
	hour := time.Now().Local().Format("15:04:05.999999")
	fmt.Println("Running: " + helper.FuncNameFromPC(pc) + ", time: " + hour)
}

func printFinishRunning(pc uintptr) {
	funcName := helper.FuncNameFromPC(pc)
	hour := time.Now().Local().Format("15:04:05.999999")
	err := recover()
	if err != nil {
		fmt.Println("Error running: " + funcName + ", time: " + hour)
		os.Exit(1)
	} else {
		fmt.Println("Finish running: " + funcName + ", time: " + hour)
	}
}
