package seeders

import (
	"context"
	"database/sql"
	"math/rand"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"time"
)

type UserRoleSpaceSeeder struct {
	db        *sql.DB
	repo      *repositories.UserRoleSpaceRepo
	userRepo  *repositories.UserRepo
	roleRepo  *repositories.RoleRepo
	spaceRepo *repositories.SpaceRepo
}

func NewUserRoleSpaceSeeder(db *sql.DB) *UserRoleSpaceSeeder {
	return &UserRoleSpaceSeeder{
		db:        db,
		repo:      repositories.NewUserRoleSpaceRepo(db),
		userRepo:  repositories.NewUserRepo(db),
		roleRepo:  repositories.NewRoleRepo(db),
		spaceRepo: repositories.NewSpaceRepo(db),
	}
}

func (seeder *UserRoleSpaceSeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	users, err := seeder.userRepo.All(ctx)
	helper.PanicIfError(err)

	roleSpaceOwner, err := seeder.roleRepo.FindByName(ctx, "space owner")
	helper.PanicIfError(err)
	roleSpaceMember, err := seeder.roleRepo.FindByName(ctx, "space member")
	helper.PanicIfError(err)

	spaces, err := seeder.spaceRepo.All(ctx)
	helper.PanicIfError(err)

	userRoleSpaces := []entities.UserRoleSpace{}

	for _, user := range users {
		role := entities.Role{}
		if helper.BooleanRandom() {
			role = roleSpaceOwner
		} else {
			role = roleSpaceMember
		}

		space := spaces[rand.Intn(len(spaces))]

		urs := entities.UserRoleSpace{
			UserId:    user.Id,
			RoleId:    role.Id,
			SpaceId:   space.Id,
			CreatedAt: time.Now(),
			UpdatedAt: helper.RandomSQLNullTime(time.Now().AddDate(0, 0, 1)),
		}

		if user.Email == "arfan@gmail.com" { // Give ownership role if email is match
			urs.RoleId = roleSpaceOwner.Id
		}

		userRoleSpaces = append(userRoleSpaces, urs)
	}

	_, err = seeder.repo.Insert(ctx, userRoleSpaces...)
	helper.PanicIfError(err)
}
