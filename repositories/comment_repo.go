package repositories

import "database/sql"

type CommentRepo struct {
	tableName string
	db        *sql.DB
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{
		tableName: "comments",
		db:        db,
	}
}
