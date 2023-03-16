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

type SpaceSeeder struct {
	db             *sql.DB
	repository     *repositories.SpaceRepository
	userRepository *repositories.UserRepository
}

func NewSpaceSeeder(db *sql.DB) *SpaceSeeder {
	return &SpaceSeeder{
		db:             db,
		repository:     repositories.NewSpaceRepository(db),
		userRepository: repositories.NewUserRepository(db),
	}
}

func (seeder *SpaceSeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	users, err := seeder.userRepository.All(ctx)
	errorh.LogPanic(err)
	var spaces []*entities.Space
	for i := 0; i < len(users); i++ {
		for i := 0; i < 2; i++ {
			fakeSpace := factories.FakeSpace()
			spaces = append(spaces, &fakeSpace)
		}
	}

	_, err = seeder.repository.Insert(ctx, spaces...)
	errorh.LogPanic(err)
}
