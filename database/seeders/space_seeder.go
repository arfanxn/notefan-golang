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
	db         *sql.DB
	repository *repositories.SpaceRepository
}

func NewSpaceSeeder(db *sql.DB) *SpaceSeeder {
	return &SpaceSeeder{
		db:         db,
		repository: repositories.NewSpaceRepository(db),
	}
}

func (seeder *SpaceSeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	totalRows := 50
	spaces := []entities.Space{}

	for i := 0; i < totalRows; i++ {
		spaces = append(spaces, factories.FakeSpace())
	}

	_, err := seeder.repository.Insert(ctx, spaces...)
	errorh.Panic(err)
}
