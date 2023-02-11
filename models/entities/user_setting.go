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
