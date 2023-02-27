package repositories

import (
	"context"
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/notefan-golang/config"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/entities"

	"github.com/google/uuid"
)

type MediaRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewMediaRepository(db *sql.DB) *MediaRepository {
	columnNames := []string{
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
	return &MediaRepository{
		db:          db,
		tableName:   "medias",
		columnNames: columnNames,
	}
}

func (repository *MediaRepository) All(ctx context.Context) ([]entities.Media, error) {
	query := "SELECT " + helper.DBSliceColumnsToStr(repository.columnNames) + " FROM " + repository.tableName
	medias := []entities.Media{}
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		helper.ErrorLog(err)
		return medias, err
	}

	for rows.Next() {
		media := entities.Media{}
		err := rows.Scan(
			&media.Id,
			&media.ModelType,
			&media.ModelId,
			&media.CollectionName,
			&media.Name,
			&media.FileName,
			&media.MimeType,
			&media.Disk,
			&media.ConversionDisk,
			&media.Size,
			&media.CreatedAt,
			&media.UpdatedAt,
		)
		if err != nil {
			helper.ErrorLog(err)
			return medias, err
		}
		medias = append(medias, media)
	}

	if len(medias) == 0 {
		return medias, exceptions.HTTPNotFound
	}

	return medias, nil
}

/* // TODO: Refactor inside of this function */
// Insert inserts medias metadata into the database
// and save the media file to the storage based on specified media disk (filesystem disk)
func (repository *MediaRepository) Insert(ctx context.Context, medias ...entities.Media) ([]entities.Media, error) {
	query := buildBatchInsertQuery(repository.tableName, len(medias), repository.columnNames...)
	valueArgs := []any{}

	for _, media := range medias {

		if media.Id == uuid.Nil {
			media.Id = uuid.New()
		}
		if media.CreatedAt.IsZero() {
			media.CreatedAt = time.Now()
		}
		if media.CollectionName == "" {
			media.CollectionName = "default"
		}
		if strings.Contains(media.FileName, "/") {
			media.FileName = filepath.Base(media.FileName)
		}
		if media.MimeType == "" {
			mmtype, err := mimetype.DetectReader(media.File)
			if err != nil {
				helper.ErrorLog(err)
				return medias, exceptions.FileInvalidType
			}
			media.MimeType = mmtype.String()
		}

		// If file exists do write operation
		root := config.FSDisks[media.Disk].Root
		path := filepath.Join(root, "medias", media.Id.String(), filepath.Base(media.FileName))

		if media.File.Len() <= 0 { // check if file exists, if not exists return an error
			return medias, exceptions.FileNotProvided
		}

		err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
		fileDst, err := os.Create(path)
		defer fileDst.Close()
		if err != nil {
			helper.ErrorLog(err)
			return medias, err
		}

		_, err = io.Copy(fileDst, media.File)
		if err != nil {
			helper.ErrorLog(err)
			return medias, err
		}

		valueArgs = append(valueArgs,
			media.Id,
			media.ModelType,
			media.ModelId,
			media.CollectionName,
			media.Name,
			media.FileName,
			media.MimeType,
			media.Disk,
			media.ConversionDisk,
			media.Size,
			media.CreatedAt,
			media.UpdatedAt,
		)
	}

	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		helper.ErrorLog(err)
		return medias, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.ErrorLog(err)
		return medias, err
	}
	return medias, nil
}

func (repository *MediaRepository) Create(ctx context.Context, media entities.Media) (entities.Media, error) {
	medias, err := repository.Insert(ctx, media)
	if err != nil {
		return entities.Media{}, err
	}

	return medias[0], nil
}
