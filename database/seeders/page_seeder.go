package seeders

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/database/factories"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/sliceh"
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
	errorh.LogPanic(err)

	var pages []*entities.Page

	for _, space := range spaces {
		for i := 0; i < 2; i++ {
			page := factories.FakePage()
			page.SpaceId = space.Id
			pages = append(pages, &page)
		}
	}

	for _, chunk := range sliceh.Chunk(pages, 100) {
		_, err = seeder.repository.Insert(ctx, chunk...)
		errorh.LogPanic(err)
	}
}
