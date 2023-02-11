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
}
