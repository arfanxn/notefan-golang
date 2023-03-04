package repositories

import (
	"bytes"
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/fileh"
	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"

	"github.com/google/uuid"
)

type MediaRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
	mutex       sync.Mutex
	waitGroup   *sync.WaitGroup
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
		mutex:       sync.Mutex{},
		waitGroup:   new(sync.WaitGroup),
		tableName:   "medias",
		columnNames: columnNames,
	}
}

/*
 * ----------------------------------------------------------------
 * Repository utilty methods ⬇
 * ----------------------------------------------------------------
 */

// scanRows scans rows of the database and returns it as structs, and returns error if any error has occurred.
func (repository *MediaRepository) scanRows(rows *sql.Rows) (medias []entities.Media, err error) {
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
		errorh.LogPanic(err) // panic if scan fails
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

/*
 * ----------------------------------------------------------------
 * Repository query methods ⬇
 * ----------------------------------------------------------------
 */

// All retrieves all data on table from database
func (repository *MediaRepository) All(ctx context.Context) ([]entities.Media, error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.columnNames) + " FROM " + repository.tableName
	rows, err := repository.db.QueryContext(ctx, query)
	errorh.LogPanic(err)
	return repository.scanRows(rows)
}

// FindByModelAndCollectionName
func (repository *MediaRepository) FindByModelAndCollectionName(
	ctx context.Context, modelTyp string, modelId string, collectionName string,
) (entities.Media, error,
) {
	queryBuf := bytes.NewBufferString("SELECT ")
	queryBuf.WriteString(stringh.SliceColumnToStr(repository.columnNames))
	queryBuf.WriteString(" FROM ")
	queryBuf.WriteString(repository.tableName)
	queryBuf.WriteString(" WHERE model_type = ? and model_id = ? and collection_name = ?")

	rows, err := repository.db.QueryContext(ctx, queryBuf.String(), modelTyp, modelId, collectionName)
	errorh.LogPanic(err)
	return repository.scanRow(rows)
}

// Insert inserts medias into the database
// and save the media file to the storage based on specified media disk (filesystem disk)
func (repository *MediaRepository) Insert(ctx context.Context, medias ...*entities.Media) (sql.Result, error) {
	query := buildBatchInsertQuery(repository.tableName, len(medias), repository.columnNames...)
	valueArgs := []any{}

	var err error
	var savedFilePaths []string

	for _, media := range medias {
		if err != nil {
			fileh.BatchRemove(savedFilePaths...) // rollback saved files
			return nil, err
		}

		repository.waitGroup.Add(1)

		go func(wg *sync.WaitGroup, media *entities.Media) {
			defer wg.Done()

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

			// check if file exists, if not exists return an error
			if media.File.Len() <= 0 {
				err = exceptions.FileNotProvided
				return
			}

			media.GuessMimeType()

			// Save media file
			err = media.SaveFile()

			repository.mutex.Lock()
			savedFilePaths = append(savedFilePaths, media.GetFilePath()) // assign media file path to "savedFilePaths" in case of error happen it will used for rollbacking the saved files
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
			repository.mutex.Unlock()

		}(repository.waitGroup, media)
	}

	repository.waitGroup.Wait()

	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		errorh.Log(err)
		fileh.BatchRemove(savedFilePaths...) // rollback saved files
		return result, err
	}

	return result, nil
}

// Create do the same thing as Insert but singularly
func (repository *MediaRepository) Create(ctx context.Context, media *entities.Media) (sql.Result, error) {
	result, err := repository.Insert(ctx, media)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (repository *MediaRepository) UpdateById(ctx context.Context, media *entities.Media) (sql.Result, error) {
	query := buildUpdateQuery(repository.tableName, repository.columnNames...) + " WHERE id = ?"

	if media.CollectionName == "" {
		media.CollectionName = "default"
	}
	if strings.Contains(media.FileName, "/") {
		media.FileName = filepath.Base(media.FileName)
	}

	// if file not provided its mean no file changes in this update
	if media.File.Len() == 0 {
		err := media.RenameFile() // rename file incase of media.Filename is updated
		if err != nil {
			errorh.Log(err)
			return nil, err
		}
	} else if media.File.Len() > 0 { // check if file is provided
		media.UpdateFile()
	}

	// Refresh entity updated at
	media.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	result, err := repository.db.ExecContext(ctx, query,
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
		media.Id)

	return result, err
}
