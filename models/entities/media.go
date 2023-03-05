package entities

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/notefan-golang/config"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/fileh"
	"github.com/notefan-golang/models/requests/file_reqs"
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
	File *file_reqs.File `json:"-"`
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
	return fileh.Save(media.GetFilePath(), media.File.Buffer)
}

// RenameFile renames the media file
func (media *Media) RenameFile() error {
	filenames, err := fileh.FileNamesFromDir(filepath.Dir(media.GetFilePath()))
	oldFile, err := os.Open(filenames[0])
	errorh.LogPanic(err)
	defer oldFile.Close()

	return os.Rename(oldFile.Name(), media.GetFilePath())
}

// UpdateFile will remove the old file and save the new file
func (media *Media) UpdateFile() error {
	err := media.RemoveDirFile()
	if err != nil {
		return err
	}
	return media.SaveFile()
}

// RemoveFile removes the media file
func (media *Media) RemoveFile() error {
	return os.Remove(media.GetFilePath())
}

// RemoveDirFile removes directory of media file
func (media *Media) RemoveDirFile() error {
	return os.RemoveAll(filepath.Dir(media.GetFilePath()))
}
