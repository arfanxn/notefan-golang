package seeders

import (
	"context"
	"database/sql"
	"math/rand"
	"time"

	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"
)

type UserRoleSpaceSeeder struct {
	db              *sql.DB
	repository      *repositories.UserRoleSpaceRepository
	userRepository  *repositories.UserRepository
	roleRepository  *repositories.RoleRepository
	spaceRepository *repositories.SpaceRepository
}

func NewUserRoleSpaceSeeder(db *sql.DB) *UserRoleSpaceSeeder {
	return &UserRoleSpaceSeeder{
		db:              db,
		repository:      repositories.NewUserRoleSpaceRepository(db),
		userRepository:  repositories.NewUserRepository(db),
		roleRepository:  repositories.NewRoleRepository(db),
		spaceRepository: repositories.NewSpaceRepository(db),
	}
}

func (seeder *UserRoleSpaceSeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	users, err := seeder.userRepository.All(ctx)
	helper.ErrorPanic(err)

	roleSpaceOwner, err := seeder.roleRepository.FindByName(ctx, "space owner")
	helper.ErrorPanic(err)
	roleSpaceMember, err := seeder.roleRepository.FindByName(ctx, "space member")
	helper.ErrorPanic(err)

	spaces, err := seeder.spaceRepository.All(ctx)
	helper.ErrorPanic(err)

	userRoleSpaces := []entities.UserRoleSpace{}

	for _, user := range users {
		role := entities.Role{}
		if helper.BoolRandom() {
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
			UpdatedAt: helper.DBRandNullOrTime(time.Now().AddDate(0, 0, 1)),
		}

		if user.Email == "arfan@gmail.com" { // Give ownership role if email is match
			urs.RoleId = roleSpaceOwner.Id
		}

		userRoleSpaces = append(userRoleSpaces, urs)
	}

	_, err = seeder.repository.Insert(ctx, userRoleSpaces...)
	helper.ErrorPanic(err)
}
