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

type PageSeeder struct {
	db        *sql.DB
	tableName string
	repo      *repositories.PageRepo
	spaceRepo *repositories.SpaceRepo
}

func NewPageSeeder(db *sql.DB) *PageSeeder {
	return &PageSeeder{
		db:        db,
		tableName: "pages",
		repo:      repositories.NewPageRepo(db),
		spaceRepo: repositories.NewSpaceRepo(db),
	}
}

func (seeder *PageSeeder) Run() {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunning(pc)
	defer printFinishRunning(pc)

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	spaces, err := seeder.spaceRepo.All(ctx)

	totalRows := len(spaces) * 2
	pages := []entities.Page{}

	for i := 0; i < totalRows; i++ {
		space := spaces[rand.Intn(len(spaces))]
		page := factories.NewPage()
		page.SpaceId = space.Id
		pages = append(pages, page)
	}

	_, err = seeder.repo.Insert(ctx, pages...)
	helper.PanicIfError(err)
}