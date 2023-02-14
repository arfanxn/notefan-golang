package seeders

import (
	"context"
	"database/sql"
	"notefan-golang/database/factories"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"time"
)

type UserSettingSeeder struct {
	db             *sql.DB
	repository     *repositories.UserSettingRepository
	userRepository *repositories.UserRepository
}

func NewUserSettingSeeder(db *sql.DB) *UserSettingSeeder {
	return &UserSettingSeeder{
		db:             db,
		repository:     repositories.NewUserSettingRepository(db),
		userRepository: repositories.NewUserRepository(db),
	}
}

func (seeder *UserSettingSeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	users, err := seeder.userRepository.All(ctx)
	helper.ErrorPanic(err)

	userSettings := []entities.UserSetting{}

	for _, user := range users {
		for i := 0; i < 5; i++ {
			userSetting := factories.FakeUserSetting()
			userSetting.UserId = user.Id
			userSettings = append(userSettings, userSetting)
		}
	}

	_, err = seeder.repository.Insert(ctx, userSettings...)
	helper.ErrorPanic(err)

}
