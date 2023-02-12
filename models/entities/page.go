package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Page struct {
	Id        uuid.UUID    `json:"id"`
	SpaceId   uuid.UUID    `json:"space_id"`
	Title     string       `json:"title"`
	Order     int          `json:"order"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
