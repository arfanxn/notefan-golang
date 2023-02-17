package repositories

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"time"

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

// scanRows scans rows of the database and returns it as structs, and returns error if any error has occurred.
func (repository *MediaRepository) scanRows(rows *sql.Rows) ([]entities.Media, error) {
	medias := []entities.Media{}
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
		helper.ErrorPanic(err) // panic if scan fails
		medias = append(medias, media)
	}

	if len(medias) == 0 {
		return medias, exceptions.HTTPNotFound
	}
	return medias, nil
}

// scanRow scans only a row of the database and returns it as struct, and returns error if any error has occurred.
func (repository *MediaRepository) scanRow(rows *sql.Rows) (entities.Media, error) {
	medias, err := repository.scanRows(rows)
	if err != nil {
		return entities.Media{}, err
	}
	return medias[0], nil
}

// All
func (repository *MediaRepository) All(ctx context.Context) ([]entities.Media, error) {
	query := "SELECT " + helper.DBSliceColumnsToStr(repository.columnNames) + " FROM " + repository.tableName
	rows, err := repository.db.QueryContext(ctx, query)
	helper.ErrorPanic(err) // panic if query error
	return repository.scanRows(rows)
}

func (repository *MediaRepository) FindByModelAndCollectionName(
	ctx context.Context, modelType string, modelId string, collectionName string,
) (entities.Media, error) {
	query := "SELECT " + helper.DBSliceColumnsToStr(repository.columnNames) + " FROM " +
		repository.tableName + " WHERE `model_type` = ? AND `model_id` = ? AND `collection_name` = ?"
	rows, err := repository.db.QueryContext(ctx, query, modelType, modelId, collectionName)
	helper.ErrorPanic(err) // panic if query error
	return repository.scanRow(rows)
}

// Insert insert medias into database
func (repository *MediaRepository) Insert(ctx context.Context, medias ...entities.Media) ([]entities.Media, error) {
	query := buildBatchInsertQuery(repository.tableName, len(medias), repository.columnNames...)
	valueArgs := []any{}

	for _, media := range medias {
		if media.Id == uuid.Nil {
			media.Id = uuid.New()
		}
		if media.FileName == "" {
			media.FileName = filepath.Base(media.File.Name())
		}
		if strings.Contains(media.FileName, "/") {
			media.FileName = filepath.Base(media.File.Name())
		}
		if media.CreatedAt.IsZero() {
			media.CreatedAt = time.Now()
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
