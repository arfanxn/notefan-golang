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

type PageContentSeeder struct {
	db        *sql.DB
	tableName string
	repo      *repositories.PageContentRepo
	pageRepo  *repositories.PageRepo
}

func NewPageContentSeeder(db *sql.DB) *PageContentSeeder {
	return &PageContentSeeder{
		db:        db,
		tableName: "page_contents",
		repo:      repositories.NewPageContentRepo(db),
		pageRepo:  repositories.NewPageRepo(db),
	}
}

func (seeder *PageContentSeeder) Run() {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunningSeeder(pc)
	defer printFinishRunningSeeder(pc)

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	pages, err := seeder.pageRepo.All(ctx)
	helper.PanicIfError(err)

	pageContents := []entities.PageContent{}

	for _, page := range pages {
		for i := 1; i <= 5; i++ {
			pageContent := factories.NewPageContent()
			pageContent.PageId = page.Id
			pageContent.Order = i

			pageContents = append(pageContents, pageContent)
		}
	}

	_, err = seeder.repo.Insert(ctx, pageContents...)
	helper.PanicIfError(err)
}
