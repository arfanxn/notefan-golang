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
	// Consoler (notify seeder has started running)
	fmt.Println("Running Seeder...")

	// ---- Begin ----
	db := seeder.db

	// Inject entity seeders into struct's field
	seeder.Seeders = append([]SeederContract{
		// User and related seeders
		NewUserSeeder(db),
		NewUserSettingSeeder(db),
		NewPermissionSeeder(db),
		NewRoleSeeder(db),
		NewPermissionRoleSeeder(db),

		// Space and related seeders
		NewSpaceSeeder(db),
		NewUserRoleSpaceSeeder(db),

		// Page and related seeders
		NewPageSeeder(db),
		NewPageContentSeeder(db),
		NewPageContentChangeHistorySeeder(db),
	}, seeder.Seeders...)

	// run seeder one by one
	for _, entitySeeder := range seeder.Seeders {
		entitySeeder.Run()
	}

	// Notify if the seeder has finished and succeeded
	printDividerLine()
	fmt.Println("Seeding completed successfully")
	os.Exit(0)
}

func printStartRunning(pc uintptr) {
	hour := time.Now().Local().Format("15:04:05.999999")
	printDividerLine()
	fmt.Println("Running: " + helper.FuncNameFromPC(pc) + ", time: " + hour)
}

func printFinishRunning(pc uintptr) {
	funcName := helper.FuncNameFromPC(pc)
	hour := time.Now().Local().Format("15:04:05.999999")
	err := recover()
	if err != nil {
		fmt.Println("Error running: " + funcName + ", time: " + hour)
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Finish running: " + funcName + ", time: " + hour)
	}
}

func printDividerLine() {
	fmt.Println("----------------------------------------------------------------")
}
