package entities

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/notefan-golang/config"
	"github.com/notefan-golang/helpers/fileh"
	"github.com/notefan-golang/models/requests/file_reqs"
)

type Media struct {
	// Id will be generated automatically, can be set manually if needed
	Id uuid.UUID `json:"id"`
	// ModelType must be specified
	ModelType string `json:"model_type"`
	// ModelId must be specified
	ModelId uuid.UUID `json:"model_id"`
	// CollectionName will be autofilled with default CollectionName if not specified
	CollectionName string `json:"collection_name"`
	// Name can be null if not specified
	Name sql.NullString `json:"name"`
	// FileName will be autofilled by random alphanumeric characters if not specified
	FileName string `json:"file_name"`
	// MimeType will be autofilled by guessing the file bytes if not specified
	MimeType string `json:"mime_type"`
	// Disk will be autofilled with default disk if not specified
	Disk string `json:"disk"`
	// ConversionDisk can be null if not set
	ConversionDisk sql.NullString `json:"conversion_disk"`
	// if not set will be autofilled with by guessing the file bytes size
	Size int64 `json:"size"`
	// CreatedAt will be autofilled on creation
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt will be autofilled after updation
	UpdatedAt sql.NullTime `json:"updated_at"`

	// File Metadata, not in table columns
	File *file_reqs.File `json:"-"`
}

/*
 * ----------------------------------------------------------------
 * Media Table and Columns methods  ⬇
 * ----------------------------------------------------------------
 */

// GetColumnNames returns the column names of the entity
func (ety Media) GetColumnNames() []string {
	return []string{
		"id",
		"model_type",
		"model_id",
		"collection_name",
		"name",
		"file_name",
		"mime_type",
		"disk",
		"conversion_disk",
		"size",
		"created_at",
		"updated_at",
	}
}

// GetTableName returns the table name
func (ety Media) GetTableName() string {
	return "medias"
}

/*
 * ----------------------------------------------------------------
 * Media File related methods ⬇
 * ----------------------------------------------------------------
 */

// GetDefaultDisk returns media default configured disk on env variable
func (media *Media) GetDefaultDisk() string {
	return os.Getenv("MEDIA_DEFAULT_DISK")
}

// GetDisk returns media disk or default disk if media disk is not set
func (media *Media) GetDisk() string {
	if media.Disk == "" {
		return media.GetDefaultDisk()
	}
	return media.Disk
}

// GetFilePath returns the path to the media's file path (media save file location)
func (media *Media) GetFilePath() string {
	return filepath.Join(
		config.FSDisks[media.GetDisk()].Root,
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
	if err != nil {
		return err
	}
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
