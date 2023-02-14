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

type CommentSeeder struct {
	db                    *sql.DB
	repository            *repositories.CommentRepository
	userRepository        *repositories.UserRepository
	pageContentRepository *repositories.PageContentRepository
}

func NewCommentSeeder(db *sql.DB) *CommentSeeder {
	return &CommentSeeder{
		db:                    db,
		repository:            repositories.NewCommentRepository(db),
		userRepository:        repositories.NewUserRepository(db),
		pageContentRepository: repositories.NewPageContentRepository(db),
	}
}

func (seeder *CommentSeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	users, err := seeder.userRepository.All(ctx)
	helper.ErrorPanic(err)

	pageContents, err := seeder.pageContentRepository.All(ctx)
	helper.ErrorPanic(err)

	comments := []entities.Comment{}

	for _, user := range users {
		for i := 0; i < 5; i++ { // each user has 5 comments
			pageContent := pageContents[rand.Intn(len(pageContents))]

			comment := factories.FakeComment()
			comment.CommentableType = helper.ReflectGetTypeName(pageContent)
			comment.CommentableId = pageContent.Id
			comment.UserId = user.Id

			comments = append(comments, comment)
		}
	}

	_, err = seeder.repository.Insert(ctx, comments...)
	helper.ErrorPanic(err)
}
