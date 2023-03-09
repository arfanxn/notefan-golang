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

type PageContentSeeder struct {
	db             *sql.DB
	repository     *repositories.PageContentRepository
	pageRepository *repositories.PageRepository
}

func NewPageContentSeeder(db *sql.DB) *PageContentSeeder {
	return &PageContentSeeder{
		db:             db,
		repository:     repositories.NewPageContentRepository(db),
		pageRepository: repositories.NewPageRepository(db),
	}
}

func (seeder *PageContentSeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	pages, err := seeder.pageRepository.All(ctx)
	errorh.Panic(err)

	var pageContents []*entities.PageContent

	for _, page := range pages {
		for i := 1; i <= 5; i++ {
			pageContent := factories.FakePageContent()
			pageContent.PageId = page.Id
			pageContent.Order = i

			pageContents = append(pageContents, &pageContent)
		}
	}

	_, err = seeder.repository.Insert(ctx, pageContents...)
	errorh.Panic(err)
}
