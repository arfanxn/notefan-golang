package seeders

import (
	"database/sql"
	"fmt"
	"notefan-golang/helper"
	"os"
	"time"
)

type DatabaseSeeder struct {
	db *sql.DB
}

func NewDatabaseSeeder(db *sql.DB) *DatabaseSeeder {
	return &DatabaseSeeder{db: db}
}

func (seeder DatabaseSeeder) Run() {
	// Consoler
	fmt.Println("Running Seeder...")
	defer func() {
		fmt.Println("Seeding completed")
		os.Exit(0)
	}()

	// Run seeders

	UserSeeder(seeder)
	PermissionSeeder(seeder)
	RoleSeeder(seeder)
	PermissionRoleSeeder(seeder)
	UserRoleSpaceSeeder(seeder)

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
