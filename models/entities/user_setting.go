package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserSetting struct {
	Id        uuid.UUID    `json:"id"`
	UserId    uuid.UUID    `json:"user_id"`
	Key       string       `json:"key"`
	Value     string       `json:"value"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

/*
 * ----------------------------------------------------------------
 * UserSetting Table and Columns methods  ⬇
 * ----------------------------------------------------------------
 */

// GetColumnNames returns the column names of the entity
func (ety UserSetting) GetColumnNames() []string {
	return []string{
		"id",
		"user_id",
		"key",
		"value",
		"created_at",
		"updated_at",
	}
}

// GetTableName returns the table name
func (ety UserSetting) GetTableName() string {
	return "user_settings"
}
