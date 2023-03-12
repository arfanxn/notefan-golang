package seeders

import (
	"context"
	"database/sql"
	"math/rand"
	"time"

	"github.com/notefan-golang/database/factories"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"
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
	errorh.LogPanic(err)

	pageContents, err := seeder.pageContentRepository.All(ctx)
	errorh.LogPanic(err)

	var comments []*entities.Comment

	for _, user := range users {
		for i := 0; i < 5; i++ { // each user has 5 comments
			pageContent := pageContents[rand.Intn(len(pageContents))]

			comment := factories.FakeComment()
			comment.CommentableType = reflecth.GetTypeName(pageContent)
			comment.CommentableId = pageContent.Id
			comment.UserId = user.Id

			comments = append(comments, &comment)
		}
	}

	_, err = seeder.repository.Insert(ctx, comments...)
	errorh.LogPanic(err)
}
