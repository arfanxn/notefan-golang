package seeders

import (
	"context"
	"database/sql"
	"math/rand"
	"time"

	"github.com/notefan-golang/database/factories"
	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"
)

type FavouriteUserSeeder struct {
	db             *sql.DB
	repository     *repositories.FavouriteUserRepository
	userRepository *repositories.UserRepository
	pageRepository *repositories.PageRepository
}

func NewFavouriteUserSeeder(db *sql.DB) *FavouriteUserSeeder {
	return &FavouriteUserSeeder{
		db:             db,
		repository:     repositories.NewFavouriteUserRepository(db),
		userRepository: repositories.NewUserRepository(db),
		pageRepository: repositories.NewPageRepository(db),
	}
}

func (seeder *FavouriteUserSeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	users, err := seeder.userRepository.All(ctx)
	helper.ErrorPanic(err)

	pages, err := seeder.pageRepository.All(ctx)
	helper.ErrorPanic(err)

	favouriteUsers := []entities.FavouriteUser{}

	for _, user := range users {
		for i := 0; i < 2; i++ { // each user has 2 favorited pages
			page := pages[rand.Intn(len(pages))]

			favouriteUser := factories.FakeFavouriteUser()
			favouriteUser.FavouriteableType = helper.ReflectGetTypeName(page)
			favouriteUser.FavouriteableId = page.Id
			favouriteUser.UserId = user.Id

			favouriteUsers = append(favouriteUsers, favouriteUser)
		}
	}

	_, err = seeder.repository.Insert(ctx, favouriteUsers...)
	helper.ErrorPanic(err)
}
