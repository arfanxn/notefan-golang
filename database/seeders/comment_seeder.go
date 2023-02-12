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
	db              *sql.DB
	repo            *repositories.CommentRepo
	userRepo        *repositories.UserRepo
	pageContentRepo *repositories.PageContentRepo
}

func NewCommentSeeder(db *sql.DB) *CommentSeeder {
	return &CommentSeeder{
		db:              db,
		repo:            repositories.NewCommentRepo(db),
		userRepo:        repositories.NewUserRepo(db),
		pageContentRepo: repositories.NewPageContentRepo(db),
	}
}

func (seeder *CommentSeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	users, err := seeder.userRepo.All(ctx)
	helper.PanicIfError(err)

	pageContents, err := seeder.pageContentRepo.All(ctx)
	helper.PanicIfError(err)

	comments := []entities.Comment{}

	for _, user := range users {
		for i := 0; i < 5; i++ { // each user has 5 comments
			pageContent := pageContents[rand.Intn(len(pageContents))]

			comment := factories.NewComment()
			comment.CommentableType = helper.GetTypeName(pageContent)
			comment.CommentableId = pageContent.Id
			comment.UserId = user.Id

			comments = append(comments, comment)
		}
	}

	_, err = seeder.repo.Insert(ctx, comments...)
	helper.PanicIfError(err)
}
