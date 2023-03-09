package seeders

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/database/factories"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"
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
	errorh.Panic(err)

	var userSettings []*entities.UserSetting

	for _, user := range users {
		for i := 0; i < 5; i++ {
			userSetting := factories.FakeUserSetting()
			userSetting.UserId = user.Id
			userSettings = append(userSettings, &userSetting)
		}
	}

	_, err = seeder.repository.Insert(ctx, userSettings...)
	errorh.Panic(err)

}
