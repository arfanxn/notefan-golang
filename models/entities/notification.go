package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	Id         uuid.UUID    `json:"id"`
	ObjectType string       `json:"object_type"`
	ObjectId   uuid.UUID    `json:"object_id"`
	Title      string       `json:"title"`
	Type       string       `json:"type"`
	Body       string       `json:"body"`
	ArchivedAt sql.NullTime `json:"archived_at"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  sql.NullTime `json:"updated_at"`
}
