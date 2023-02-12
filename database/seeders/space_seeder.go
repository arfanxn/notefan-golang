package seeders

import (
	"context"
	"database/sql"
	"notefan-golang/database/factories"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"runtime"
	"time"
)

type SpaceSeeder struct {
	db        *sql.DB
	tableName string
	repo      *repositories.SpaceRepo
}

func NewSpaceSeeder(db *sql.DB) *SpaceSeeder {
	return &SpaceSeeder{
		db:        db,
		tableName: "spaces",
		repo:      repositories.NewSpaceRepo(db),
	}
}

func (seeder *SpaceSeeder) Run() {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunningSeeder(pc)
	defer printFinishRunningSeeder(pc)

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	totalRows := 50
	spaces := []entities.Space{}

	for i := 0; i < totalRows; i++ {
		spaces = append(spaces, factories.NewSpace())
	}

	_, err := seeder.repo.Insert(ctx, spaces...)
	helper.PanicIfError(err)
}
