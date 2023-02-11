package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type FavouriteUser struct {
	Id                uuid.UUID    `json:"id"`
	FavouriteableType string       `json:"favouriteable_type"`
	FavouriteableId   uuid.UUID    `json:"favouriteable_id"`
	UserId            uuid.UUID    `json:"user_id"`
	Order             int          `json:"order"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         sql.NullTime `json:"updated_at"`
}
