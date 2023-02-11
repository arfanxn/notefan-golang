package entities

import "github.com/google/uuid"

type Permission struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
