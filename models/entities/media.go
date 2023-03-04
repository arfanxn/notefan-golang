package entities

import (
	"bytes"
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/notefan-golang/config"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/fileh"
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
		mime, err := mimetype.DetectReader(media.File)
		if err != nil {
			errorh.LogPanic(err)
		}
		media.MimeType = mime.String()
	}
}

// GetFilePath returns the path to the media's file path (media save file location)
func (media *Media) GetFilePath() string {
	return filepath.Join(
		config.FSDisks[media.Disk].Root,
		"medias",
		media.Id.String(),
		filepath.Base(media.FileName),
	)
}

// SaveFile saves the media file
func (media *Media) SaveFile() error {
	return fileh.Save(media.GetFilePath(), media.File)
}

// RenameFile renames the media file
func (media *Media) RenameFile() error {

	// Get random file from directory, the random result will be the same since the directory it self only contains a single file that belongs to the media
	oldFile, err := fileh.RandFromDir(filepath.Dir(media.GetFilePath()))
	errorh.LogPanic(err)
	defer oldFile.Close()

	return os.Rename(oldFile.Name(), media.GetFilePath())
}

// RemoveFile removes the media file
func (media *Media) RemoveFile() error {
	return os.Remove(media.GetFilePath())
}

// RemoveDirFile removes directory of media file
func (media *Media) RemoveDirFile() error {
	return os.RemoveAll(filepath.Dir(media.GetFilePath()))
}
