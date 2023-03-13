package entities

import "github.com/google/uuid"

type Permission struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

/*
 * ----------------------------------------------------------------
 * Permission Table and Columns methods  ⬇
 * ----------------------------------------------------------------
 */

// GetColumnNames returns the column names of the entity
func (ety Permission) GetColumnNames() []string {
	return []string{
		"id",
		"name",
	}
}

// GetTableName returns the table name
func (ety Permission) GetTableName() string {
	return "permissions"
}
