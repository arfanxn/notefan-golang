package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id              uuid.UUID    `json:"id"`
	CommentableType string       `json:"commentable_type"`
	CommentableId   uuid.UUID    `json:"commentable_id"`
	UserId          uuid.UUID    `json:"user_id"`
	Body            string       `json:"body"`
	ResolvedAt      sql.NullTime `json:"resolved_at"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       sql.NullTime `json:"updated_at"`
}
