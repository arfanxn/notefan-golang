package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type CommentReaction struct {
	Id        uuid.UUID    `json:"id"`
	CommentId uuid.UUID    `json:"comment_id"`
	UserId    uuid.UUID    `json:"user_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

// GetColumnNames returns the column names of the entity
func (ety CommentReaction) GetColumnNames() []string {
	return []string{
		"id",
		"comment_id",
		"user_id",
		"created_at",
		"updated_at",
	}
}

// GetTableName returns the table name
func (ety CommentReaction) GetTableName() string {
	return "comment_reactions"
}
