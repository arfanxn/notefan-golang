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

func (this *DatabaseSeeder) Run() {
	this.notifySeederStarted()

	// ---- Begin ----
	db := this.db

	// Inject entity seeders into struct's field
	this.Seeders = append([]SeederContract{
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
	}, this.Seeders...)

	// run seeder one by one
	for _, seeder := range this.Seeders {
		func() {
			this.notifyEntitySeeederStarted(seeder)
			defer this.notifyEntitySeederFinished(seeder)
			seeder.Run()
		}()
	}

	this.notifySeederFinished()
}

func (DatabaseSeeder) notifySeederStarted() {
	// Consoler (notify seeder has started running)
	fmt.Println("")
	fmt.Println("Running Seeder...")
}

func (DatabaseSeeder) notifySeederFinished() {
	// Notify if the seeder has finished and succeeded
	fmt.Println("----------------------------------------------------------------")
	fmt.Println("Seeding completed successfully")
	os.Exit(0)
}

func (DatabaseSeeder) notifyEntitySeeederStarted(seeder SeederContract) {
	hour := time.Now().Local().Format("15:04:05.999999")
	seederName := strings.ReplaceAll(helper.ReflectGetTypeName(seeder), "*", "")
	fmt.Println("----------------------------------------------------------------")
	fmt.Println("Running: " + seederName + ", time: " + hour)
}

func (DatabaseSeeder) notifyEntitySeederFinished(seeder SeederContract) {
	hour := time.Now().Local().Format("15:04:05.999999")
	seederName := strings.ReplaceAll(helper.ReflectGetTypeName(seeder), "*", "")
	err := recover()
	if err != nil {
		fmt.Println("Error running: " + seederName + ", time: " + hour)
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Finish running: " + seederName + ", time: " + hour)
	}
}
