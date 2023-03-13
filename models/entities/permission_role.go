package entities

import (
	"time"

	"github.com/google/uuid"
)

type PermissionRole struct {
	PermissionId uuid.UUID `json:"permission_id"`
	RoleId       uuid.UUID `json:"role_id"`
	CreatedAt    time.Time `json:"created_at"`
}

/*
 * ----------------------------------------------------------------
 * PermissionRole Table and Columns methods  ⬇
 * ----------------------------------------------------------------
 */

// GetColumnNames returns the column names of the entity
func (ety PermissionRole) GetColumnNames() []string {
	return []string{
		"permission_id",
		"role_id",
		"created_at",
	}
}

// GetTableName returns the table name
func (ety PermissionRole) GetTableName() string {
	return "permission_role"
}
