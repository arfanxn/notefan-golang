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
	db                *sql.DB
	repository        *repositories.CommentReactionRepository
	commentRepository *repositories.CommentRepository
	userRepository    *repositories.UserRepository
}

func NewCommentReactionSeeder(db *sql.DB) *CommentReactionSeeder {
	return &CommentReactionSeeder{
		db:                db,
		repository:        repositories.NewCommentReactionRepository(db),
		commentRepository: repositories.NewCommentRepository(db),
		userRepository:    repositories.NewUserRepository(db),
	}
}

func (seeder *CommentReactionSeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	comments, err := seeder.commentRepository.All(ctx)
	helper.ErrorPanic(err)

	users, err := seeder.userRepository.All(ctx)
	helper.ErrorPanic(err)

	commentReactions := []entities.CommentReaction{}

	for _, comment := range comments {
		for i := 0; i < 2; i++ {
			user := users[rand.Intn(len(users))]

			commentReaction := factories.FakeCommentReaction()
			commentReaction.CommentId = comment.Id
			commentReaction.UserId = user.Id
			commentReactions = append(commentReactions, commentReaction)
		}
	}

	_, err = seeder.repository.Insert(ctx, commentReactions...)
	helper.ErrorPanic(err)
}
