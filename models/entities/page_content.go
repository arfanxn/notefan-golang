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

/*
 * ----------------------------------------------------------------
 * PageContent Table and Columns methods  ⬇
 * ----------------------------------------------------------------
 */

// GetColumnNames returns the column names of the entity
func (ety PageContent) GetColumnNames() []string {
	return []string{
		"id",
		"page_id",
		"type",
		"order",
		"body",
		"created_at",
		"updated_at",
	}
}

// GetTableName returns the table name
func (ety PageContent) GetTableName() string {
	return "page_contents"
}
