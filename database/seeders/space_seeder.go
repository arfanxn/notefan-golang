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

type SpaceSeeder struct {
	db   *sql.DB
	repo *repositories.SpaceRepo
}

func NewSpaceSeeder(db *sql.DB) *SpaceSeeder {
	return &SpaceSeeder{
		db:   db,
		repo: repositories.NewSpaceRepo(db),
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

	_, err := seeder.repo.Insert(ctx, spaces...)
	helper.ErrorPanic(err)
}
