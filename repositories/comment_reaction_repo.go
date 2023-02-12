package repositories

import "database/sql"

type CommentReactionRepo struct {
	tableName string
	db        *sql.DB
}

func NewCommentReactionRepo(db *sql.DB) *CommentReactionRepo {
	return &CommentReactionRepo{
		tableName: "comment_reactions",
		db:        db,
	}
}
