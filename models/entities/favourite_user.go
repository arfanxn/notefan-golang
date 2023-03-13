package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type FavouriteUser struct {
	FavouriteableType string       `json:"favouriteable_type"`
	FavouriteableId   uuid.UUID    `json:"favouriteable_id"`
	UserId            uuid.UUID    `json:"user_id"`
	Order             int          `json:"order"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         sql.NullTime `json:"updated_at"`
}

// GetColumnNames returns the column names of the entity
func (ety FavouriteUser) GetColumnNames() []string {
	return []string{
		"favouriteable_type",
		"favouriteable_id",
		"user_id",
		"order",
		"created_at",
		"updated_at",
	}
}

// GetTableName returns the table name
func (ety FavouriteUser) GetTableName() string {
	return "favorite_user"
}
