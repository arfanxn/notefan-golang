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

type CommentReactionSeeder struct {
	db          *sql.DB
	repo        *repositories.CommentReactionRepo
	commentRepo *repositories.CommentRepo
	userRepo    *repositories.UserRepo
}

func NewCommentReactionSeeder(db *sql.DB) *CommentReactionSeeder {
	return &CommentReactionSeeder{
		db:          db,
		repo:        repositories.NewCommentReactionRepo(db),
		commentRepo: repositories.NewCommentRepo(db),
		userRepo:    repositories.NewUserRepo(db),
	}
}

func (seeder *CommentReactionSeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	comments, err := seeder.commentRepo.All(ctx)
	helper.PanicIfError(err)

	users, err := seeder.userRepo.All(ctx)
	helper.PanicIfError(err)

	commentReactions := []entities.CommentReaction{}

	for _, comment := range comments {
		for i := 0; i < 2; i++ {
			user := users[rand.Intn(len(users))]

			commentReaction := factories.NewCommentReaction()
			commentReaction.CommentId = comment.Id
			commentReaction.UserId = user.Id
			commentReactions = append(commentReactions, commentReaction)
		}
	}

	_, err = seeder.repo.Insert(ctx, commentReactions...)
	helper.PanicIfError(err)
}
