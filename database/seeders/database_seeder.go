package seeders

import (
	"database/sql"
	"fmt"
	"notefan-golang/helper"
	"os"
	"strings"
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
	printStartRunningDBSeeder()

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

		// Notification and related seeders
		NewNotificationSeeder(db),
		NewNotificationUserSeeder(db),

		// Comment and related seeders
		NewCommentSeeder(db),
		NewCommentReactionSeeder(db),

		// Favourite or Bookmark seeder
		NewFavouriteUserSeeder(db),

		// Media Seeder
		NewMediaSeeder(db),
	}, seeder.Seeders...)

	// run seeder one by one
	for _, entitySeeder := range seeder.Seeders {
		printStartRunningEntitySeeder(seeder)
		entitySeeder.Run()
		printFinishRunningEntitySeeder(seeder)
	}

	printFinishRunningDBSeeder()
}

func printStartRunningDBSeeder() {
	// Consoler (notify seeder has started running)
	fmt.Println("")
	fmt.Println("Running Seeder...")
}

func printFinishRunningDBSeeder() {
	// Notify if the seeder has finished and succeeded
	printDividerLine()
	fmt.Println("Seeding completed successfully")
	os.Exit(0)
}

func printStartRunningEntitySeeder(seeder SeederContract) {
	hour := time.Now().Local().Format("15:04:05.999999")
	seederName := strings.ReplaceAll(helper.GetTypeName(seeder), "*", "")
	printDividerLine()
	fmt.Println("Running: " + seederName + ", time: " + hour)
}

func printFinishRunningEntitySeeder(seeder SeederContract) {
	hour := time.Now().Local().Format("15:04:05.999999")
	seederName := strings.ReplaceAll(helper.GetTypeName(seeder), "*", "")
	err := recover()
	if err != nil {
		fmt.Println("Error running: " + seederName + ", time: " + hour)
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Finish running: " + seederName + ", time: " + hour)
	}
}

func printDividerLine() {
	fmt.Println("----------------------------------------------------------------")
}
