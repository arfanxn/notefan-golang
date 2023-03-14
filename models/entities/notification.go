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

	// Relations

	Object any `json:"object"`
}

/*
 * ----------------------------------------------------------------
 * Notification Table and Columns methods  ⬇
 * ----------------------------------------------------------------
 */

// GetColumnNames returns the column names of the entity
func (ety Notification) GetColumnNames() []string {
	return []string{
		"id",
		"object_type",
		"object_id",
		"title",
		"type",
		"body",
		"archived_at",
		"created_at",
		"updated_at",
	}
}

// GetTableName returns the table name
func (ety Notification) GetTableName() string {
	return "notifications"
}
