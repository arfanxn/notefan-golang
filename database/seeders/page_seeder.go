package seeders

import (
	"context"
	"database/sql"
	"math/rand"
	"time"

	"github.com/notefan-golang/database/factories"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"
)

type PageSeeder struct {
	db              *sql.DB
	repository      *repositories.PageRepository
	spaceRepository *repositories.SpaceRepository
}

func NewPageSeeder(db *sql.DB) *PageSeeder {
	return &PageSeeder{
		db:              db,
		repository:      repositories.NewPageRepository(db),
		spaceRepository: repositories.NewSpaceRepository(db),
	}
}

func (seeder *PageSeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	spaces, err := seeder.spaceRepository.All(ctx)

	totalRows := len(spaces) * 2
	var pages []*entities.Page

	for i := 0; i < totalRows; i++ {
		space := spaces[rand.Intn(len(spaces))]
		page := factories.FakePage()
		page.SpaceId = space.Id
		pages = append(pages, &page)
	}

	_, err = seeder.repository.Insert(ctx, pages...)
	errorh.LogPanic(err)
}
