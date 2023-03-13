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

/*
 * ----------------------------------------------------------------
 * Page Table and Columns methods  ⬇
 * ----------------------------------------------------------------
 */

// GetColumnNames returns the column names of the entity
func (ety Page) GetColumnNames() []string {
	return []string{
		"id",
		"space_id",
		"title",
		"order",
		"created_at",
		"updated_at",
	}
}

// GetTableName returns the table name
func (ety Page) GetTableName() string {
	return "pages"
}
