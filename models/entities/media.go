package entities

import (
	"bytes"
	"database/sql"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/notefan-golang/helpers/errorh"
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

// GuessMimeType will guess the mime type by looking up the media file
func (media *Media) GuessMimeType() {
	if media.MimeType == "" {
		mmtype, err := mimetype.DetectReader(media.File)
		if err != nil {
			errorh.LogPanic(err)
		}
		media.MimeType = mmtype.String()
	}
}
