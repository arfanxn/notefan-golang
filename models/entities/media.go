package entities

import (
	"bytes"
	"database/sql"
	"io"
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

// FillFromModel fills model related fields
func (media *Media) FillFromModel(modelTyp string, modelId string) {
	media.ModelType = modelTyp
	media.ModelId = uuid.MustParse(modelId)
}

// FillFromOSFile fills file related fields from the given osFile argument
func (media *Media) FillFromOSFile(osFile *os.File) error {
	if media.File == nil {
		media.File = bytes.NewBuffer(nil)
	}

	media.File.Reset()
	_, err := io.Copy(media.File, osFile)
	if err != nil {
		errorh.Log(err)
		return err
	}

	media.FileName = filepath.Base(osFile.Name())
	media.Size = fileh.GetSize(osFile)
	media.AutofillMimeType()

	return nil
}

// AutofillMimeType will autofill media.MimeType by looking up the media file
func (media *Media) AutofillMimeType() {
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
	file := *media.File
	return fileh.Save(media.GetFilePath(), &file)
}

// RenameFile renames the media file
func (media *Media) RenameFile() error {

	// Get random file from directory, the random result will be the same since the directory it self only contains a single file that belongs to the media
	oldFile, err := fileh.RandFromDir(filepath.Dir(media.GetFilePath()))
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
