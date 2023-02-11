package seeders

import (
	"database/sql"
	"fmt"
	"notefan-golang/helper"
	"os"
	"time"
)

type Seeder struct {
	db *sql.DB
}

func NewSeeder(db *sql.DB) *Seeder {
	return &Seeder{db: db}
}

func (seeder *Seeder) Run() {
	fmt.Println("Running Seeder...")
	UserSeeder(seeder)

	fmt.Println("Seeding completed")
	os.Exit(1)
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
