package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Media struct {
	Id             uuid.UUID    `json:"id"`
	ModelType      string       `json:"model_type"`
	ModelId        uuid.UUID    `json:"model_id"`
	CollectionName string       `json:"collection_name"`
	Name           string       `json:"name"`
	FileName       string       `json:"file_name"`
	MimeType       string       `json:"mime_type"`
	Disk           string       `json:"disk"`
	ConversionDisk string       `json:"conversion_disk"`
	Size           int          `json:"size"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      sql.NullTime `json:"updated_at"`
}
