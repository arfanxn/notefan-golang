package entities

import "github.com/google/uuid"

type Role struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	// Relations

	Permissions []Permission `json:"permissions"`
}

/*
 * ----------------------------------------------------------------
 * Role Table and Columns methods  ⬇
 * ----------------------------------------------------------------
 */

// GetColumnNames returns the column names of the entity
func (ety Role) GetColumnNames() []string {
	return []string{
		"id",
		"name",
	}
}

// GetTableName returns the table name
func (ety Role) GetTableName() string {
	return "roles"
}
