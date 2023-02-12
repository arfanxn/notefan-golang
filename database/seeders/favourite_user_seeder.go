package seeders

import (
	"context"
	"database/sql"
	"math/rand"
	"notefan-golang/database/factories"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"runtime"
	"time"
)

type FavouriteUserSeeder struct {
	db        *sql.DB
	tableName string
	repo      *repositories.FavouriteUserRepo
	userRepo  *repositories.UserRepo
	pageRepo  *repositories.PageRepo
}

func NewFavouriteUserSeeder(db *sql.DB) *FavouriteUserSeeder {
	return &FavouriteUserSeeder{
		db:        db,
		tableName: "favourite_user",
		repo:      repositories.NewFavouriteUserRepo(db),
		userRepo:  repositories.NewUserRepo(db),
		pageRepo:  repositories.NewPageRepo(db),
	}
}

func (seeder *FavouriteUserSeeder) Run() {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunningSeeder(pc)
	defer printFinishRunningSeeder(pc)

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	users, err := seeder.userRepo.All(ctx)
	helper.PanicIfError(err)

	pages, err := seeder.pageRepo.All(ctx)
	helper.PanicIfError(err)

	favouriteUsers := []entities.FavouriteUser{}

	for _, user := range users {
		for i := 0; i < 2; i++ { // each user has 2 favorited pages
			page := pages[rand.Intn(len(pages))]

			favouriteUser := factories.NewFavouriteUser()
			favouriteUser.FavouriteableType = helper.GetTypeName(page)
			favouriteUser.FavouriteableId = page.Id
			favouriteUser.UserId = user.Id

			favouriteUsers = append(favouriteUsers, favouriteUser)
		}
	}

	_, err = seeder.repo.Insert(ctx, favouriteUsers...)
	helper.PanicIfError(err)
}
