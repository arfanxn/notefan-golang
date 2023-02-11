package entities

import "github.com/google/uuid"

type PermissionRole struct {
	PermissionId uuid.UUID `json:"permission_id"`
	RoleId       string    `json:"name"`
}
