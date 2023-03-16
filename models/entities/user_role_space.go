package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserRoleSpace struct {
	Id        uuid.UUID    `json:"id"`
	UserId    uuid.UUID    `json:"user_id"`
	RoleId    uuid.UUID    `json:"role_id"`
	SpaceId   uuid.UUID    `json:"space_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`

	// Relations

	User  User  `json:"user"`
	Role  Role  `json:"role"`
	Space Space `json:"space"`
}

/*
 * ----------------------------------------------------------------
 * UserRoleSpace Table and Columns methods  ⬇
 * ----------------------------------------------------------------
 */

// GetColumnNames returns the column names of the entity
func (ety UserRoleSpace) GetColumnNames() []string {
	return []string{
		"id",
		"user_id",
		"role_id",
		"space_id",
		"created_at",
		"updated_at",
	}
}

// GetTableName returns the table name
func (ety UserRoleSpace) GetTableName() string {
	return "user_role_space"
}
