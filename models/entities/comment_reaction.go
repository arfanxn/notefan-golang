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
