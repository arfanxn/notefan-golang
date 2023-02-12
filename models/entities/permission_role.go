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
