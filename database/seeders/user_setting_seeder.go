package seeders

import (
	"context"
	"database/sql"
	"notefan-golang/database/factories"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"runtime"
	"time"
)

type UserSettingSeeder struct {
	db        *sql.DB
	tableName string
	repo      *repositories.UserSettingRepo
	userRepo  *repositories.UserRepo
}

func NewUserSettingSeeder(db *sql.DB) *UserSettingSeeder {
	return &UserSettingSeeder{
		db:        db,
		tableName: "user_settings",
		repo:      repositories.NewUserSettingRepo(db),
		userRepo:  repositories.NewUserRepo(db),
	}
}

func (seeder *UserSettingSeeder) Run() {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunning(pc)
	defer printFinishRunning(pc)

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	users, err := seeder.userRepo.All(ctx)
	helper.PanicIfError(err)

	userSettings := []entities.UserSetting{}

	for _, user := range users {
		for i := 0; i < 5; i++ {
			userSetting := factories.NewUserSetting()
			userSetting.UserId = user.Id
			userSettings = append(userSettings, userSetting)
		}
	}

	_, err = seeder.repo.Insert(ctx, userSettings...)
	helper.PanicIfError(err)

}