package seeders

import (
	"context"
	"database/sql"
	"math/rand"
	"notefan-golang/database/factories"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"time"
)

type FavouriteUserSeeder struct {
	db       *sql.DB
	repo     *repositories.FavouriteUserRepo
	userRepo *repositories.UserRepo
	pageRepo *repositories.PageRepo
}

func NewFavouriteUserSeeder(db *sql.DB) *FavouriteUserSeeder {
	return &FavouriteUserSeeder{
		db:       db,
		repo:     repositories.NewFavouriteUserRepo(db),
		userRepo: repositories.NewUserRepo(db),
		pageRepo: repositories.NewPageRepo(db),
	}
}

func (seeder *FavouriteUserSeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	users, err := seeder.userRepo.All(ctx)
	helper.ErrorPanic(err)

	pages, err := seeder.pageRepo.All(ctx)
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

	_, err = seeder.repo.Insert(ctx, favouriteUsers...)
	helper.ErrorPanic(err)
}
