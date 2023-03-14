package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Space struct {
	Id          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Domain      string       `json:"domain"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   sql.NullTime `json:"updated_at"`

	// Relations

	Owner   User   `json:"owner"`
	Members []User `json:"members"`
}

// GetColumnNames returns the column names of the entity
func (ety Space) GetColumnNames() []string {
	return []string{
		"id",
		"name",
		"description",
		"domain",
		"created_at",
		"updated_at",
	}
}

// GetTableName returns the table name
func (ety Space) GetTableName() string {
	return "spaces"
}
