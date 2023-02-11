package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserRoleSpace struct {
	UserId    uuid.UUID    `json:"user_id"`
	RoleId    uuid.UUID    `json:"role_id"`
	SpaceId   uuid.UUID    `json:"space_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
