package seeders

import (
	"context"
	"database/sql"
	"math/rand"
	"notefan-golang/database/factories"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"time"
)

type PageContentChangeHistorySeeder struct {
	db                    *sql.DB
	repository            *repositories.PageContentChangeHistoryRepository
	userRepository        *repositories.UserRepository
	pageContentRepository *repositories.PageContentRepository
}

func NewPageContentChangeHistorySeeder(db *sql.DB) *PageContentChangeHistorySeeder {
	return &PageContentChangeHistorySeeder{
		db:                    db,
		repository:            repositories.NewPageContentChangeHistoryRepository(db),
		userRepository:        repositories.NewUserRepository(db),
		pageContentRepository: repositories.NewPageContentRepository(db),
	}
}

func (seeder *PageContentChangeHistorySeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	users, err := seeder.userRepository.All(ctx)
	helper.ErrorPanic(err)

	oldPageContents, err := seeder.pageContentRepository.All(ctx)
	helper.ErrorPanic(err)

	pageContentChangeHistories := []entities.PageContentChangeHistory{}

	for _, oldPageContent := range oldPageContents {
		user := users[rand.Intn(len(users))]

		newPageContents := helper.SliceFilter(oldPageContents, func(pageContent entities.PageContent) bool {
			// return only if the page content id are not equals to the old page content id
			// and both of page content has same page id
			return (pageContent.Id.String() != oldPageContent.Id.String()) &&
				(pageContent.PageId.String() == oldPageContent.PageId.String())
		})
		if len(newPageContents) == 0 { // continue if no one returned from slice filter
			continue
		}
		newPageContent := newPageContents[rand.Intn(len(newPageContents))]

		pcch := factories.FakePageContentChangeHistory()
		pcch.BeforePageContentId = oldPageContent.Id
		pcch.AfterPageContentId = newPageContent.Id
		pcch.UserId = user.Id

		// Append the page content change history
		pageContentChangeHistories = append(pageContentChangeHistories, pcch)
	}

	// Insert the page content change history into database one by one with for loop to prevent duplicate values
	for _, pcch := range pageContentChangeHistories {
		_, err = seeder.repository.Insert(ctx, pcch)
		helper.ErrorPanic(err)
	}
}
