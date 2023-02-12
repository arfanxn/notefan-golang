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
}