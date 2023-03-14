package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID    `json:"id"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	Password  string       `json:"password"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`

	// Relations

	Role   Role    `json:"role"`
	Spaces []Space `json:"spaces"`
}

/*
 * ----------------------------------------------------------------
 * User Table and Columns methods  ⬇
 * ----------------------------------------------------------------
 */

// GetColumnNames returns the column names of the entity
func (ety User) GetColumnNames() []string {
	return []string{
		"id",
		"name",
		"email",
		"password",
		"created_at",
		"updated_at",
	}
}

// GetTableName returns the table name
func (ety User) GetTableName() string {
	return "users"
}
