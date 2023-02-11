package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type PageContent struct {
	Id        uuid.UUID    `json:"id"`
	PageId    uuid.UUID    `json:"page_id"`
	Type      string       `json:"type"`
	Order     int          `json:"order"`
	Body      string       `json:"body"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
