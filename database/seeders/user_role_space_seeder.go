package seeders

import (
	"context"
	"database/sql"
	"math/rand"
	"time"

	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/nullh"
	"github.com/notefan-golang/helpers/sliceh"
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
	errorh.Panic(err)

	roleSpaceOwner, err := seeder.roleRepository.FindByName(ctx, "space owner")
	errorh.Panic(err)
	roleSpaceMember, err := seeder.roleRepository.FindByName(ctx, "space member")
	errorh.Panic(err)

	spaces, err := seeder.spaceRepository.All(ctx)
	errorh.Panic(err)

	var userRoleSpaces []*entities.UserRoleSpace

	for _, user := range users {
		spacesLen := len(spaces)
		numberOfSpaces := 2
		spacesRandStartIndex := rand.Intn(spacesLen - numberOfSpaces)
		spaces := spaces[spacesRandStartIndex : spacesRandStartIndex+numberOfSpaces]
		for _, space := range spaces {
			role := sliceh.Random([]entities.Role{roleSpaceOwner, roleSpaceMember}) // random role

			urs := entities.UserRoleSpace{
				UserId:    user.Id,
				RoleId:    role.Id,
				SpaceId:   space.Id,
				CreatedAt: time.Now(),
				UpdatedAt: nullh.RandSqlNullTime(time.Now().AddDate(0, 0, 1)),
			}

			if user.Email == "arfan@gmail.com" { // Give ownership role if email is match
				urs.RoleId = roleSpaceOwner.Id
			}

			userRoleSpaces = append(userRoleSpaces, &urs)
		}
	}

	_, err = seeder.repository.Insert(ctx, userRoleSpaces...)
	errorh.Panic(err)
}
