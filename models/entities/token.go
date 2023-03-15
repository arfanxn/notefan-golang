package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	Id            uuid.UUID    `json:"id"`
	TokenableType string       `json:"tokenable_string"`
	TokenableId   uuid.UUID    `json:"tokenable_id"`
	Type          string       `json:"type"`
	Body          string       `json:"body"` // the token content/body/string
	UsedAt        sql.NullTime `json:"used_at"`
	ExpiredAt     sql.NullTime `json:"expired_at"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     sql.NullTime `json:"updated_at"`

	// Relations

	Tokenable any `json:"tokenable"`
}

/*
 * ----------------------------------------------------------------
 * Token Table and Columns methods  ⬇
 * ----------------------------------------------------------------
 */

// GetColumnNames returns the column names of the entity
func (ety Token) GetColumnNames() []string {
	return []string{
		"id",
		"tokenable_type",
		"tokenable_id",
		"type",
		"body",
		"used_at",
		"expired_at",
		"created_at",
		"updated_at",
	}
}

// GetTableName returns the table name
func (ety Token) GetTableName() string {
	return "tokens"
}

/*
 * ----------------------------------------------------------------
 * Token methods  ⬇
 * ----------------------------------------------------------------
 */

// IsExpired returns bool that determines the entity is expired or not
func (ety Token) IsExpired() bool {
	return ety.ExpiredAt.Time.Before(time.Now()) // check whether entity expired
}
