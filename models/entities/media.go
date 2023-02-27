package entities

import (
	"bytes"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Media struct {
	Id             uuid.UUID      `json:"id"`
	ModelType      string         `json:"model_type"`
	ModelId        uuid.UUID      `json:"model_id"`
	CollectionName string         `json:"collection_name"`
	Name           sql.NullString `json:"name"`
	FileName       string         `json:"file_name"`
	MimeType       string         `json:"mime_type"`
	Disk           string         `json:"disk"`
	ConversionDisk sql.NullString `json:"conversion_disk"`
	Size           int64          `json:"size"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      sql.NullTime   `json:"updated_at"`

	// Metadata
	File *bytes.Buffer `json:"-"`
}
