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
	"github.com/notefan-golang/helpers/fileh"
	"github.com/notefan-golang/helpers/sliceh"
	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/helpers/synch"
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
func (repository *MediaRepository) scanRows(rows *sql.Rows) (
	medias []entities.Media, err error) {
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
			return medias, err
		}
		medias = append(medias, media)
	}

	return medias, nil
}

// scanRow scans only a row of the database and returns it as struct, and returns error if any error has occurred.
func (repository *MediaRepository) scanRow(rows *sql.Rows) (media entities.Media, err error) {
	medias, err := repository.scanRows(rows)
	if err != nil {
		return
	}
	if len(medias) == 0 {
		return media, nil
	}
	media = medias[0]
	return media, nil
}

/*
 * ----------------------------------------------------------------
 * Repository query methods ⬇
 * ----------------------------------------------------------------
 */

// All retrieves all data on table from database
func (repository *MediaRepository) All(ctx context.Context) (
	medias []entities.Media, err error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.columnNames) + " FROM " + repository.tableName
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		return
	}
	medias, err = repository.scanRows(rows)
	if err != nil {
		return
	}
	return
}

// Find retrieves data on table from database by the given id
func (repository *MediaRepository) Find(ctx context.Context, id string) (media entities.Media, err error) {
	queryBuf := bytes.NewBufferString("SELECT ")
	queryBuf.WriteString(stringh.SliceColumnToStr(repository.columnNames))
	queryBuf.WriteString(" FROM ")
	queryBuf.WriteString(repository.tableName)
	queryBuf.WriteString(" WHERE `id` = ?")
	rows, err := repository.db.QueryContext(ctx, queryBuf.String(), id)
	if err != nil {
		return
	}
	media, err = repository.scanRow(rows)
	if err != nil {
		return
	}
	return
}

// FindByModelAndCollectionName
func (repository *MediaRepository) FindByModelAndCollectionName(
	ctx context.Context, modelTyp string, modelId string, collectionName string,
) (media entities.Media, err error,
) {
	queryBuf := bytes.NewBufferString("SELECT ")
	queryBuf.WriteString(stringh.SliceColumnToStr(repository.columnNames))
	queryBuf.WriteString(" FROM ")
	queryBuf.WriteString(repository.tableName)
	queryBuf.WriteString(" WHERE model_type = ? AND model_id = ? AND collection_name = ?")

	rows, err := repository.db.QueryContext(ctx, queryBuf.String(), modelTyp, modelId, collectionName)
	if err != nil {
		return
	}
	media, err = repository.scanRow(rows)
	if err != nil {
		return
	}
	return
}

// GetByModelsAndCollectionNames get data on table from database by model_types, model_ids. collection_names
func (repository *MediaRepository) GetByModelsAndCollectionNames(ctx context.Context, medias ...entities.Media) (
	[]entities.Media, error) {
	var valueArgs []any
	queryBuf := bytes.NewBufferString("SELECT ")
	queryBuf.WriteString(stringh.SliceColumnToStr(repository.columnNames))
	queryBuf.WriteString(" FROM ")
	queryBuf.WriteString(repository.tableName)

	queryBuf.WriteString(" WHERE ")
	for index, media := range medias {
		if index != 0 {
			queryBuf.WriteString(" OR ")
		}
		queryBuf.WriteString("(`model_type` = ? AND `model_id` = ? AND `collection_name` = ?)")

		valueArgs = append(valueArgs,
			media.ModelType,
			media.ModelId,
			media.CollectionName,
		)
	}

	rows, err := repository.db.QueryContext(ctx, queryBuf.String(), valueArgs...)
	if err != nil {
		return []entities.Media{}, err
	}
	return repository.scanRows(rows)
}

// GetByIds get data on table from database by the given ids
func (repository *MediaRepository) GetByIds(ctx context.Context, ids ...string) (medias []entities.Media, err error) {
	queryBuf := bytes.NewBufferString("SELECT ")
	queryBuf.WriteString(stringh.SliceColumnToStr(repository.columnNames))
	queryBuf.WriteString(" FROM ")
	queryBuf.WriteString(repository.tableName)
	queryBuf.WriteString(" WHERE `id` IN (?" + strings.Repeat(", ?", len(ids)-1) + ")")
	valueArgs := sliceh.Map(ids, func(id string) any {
		return any(id)
	})
	rows, err := repository.db.QueryContext(ctx, queryBuf.String(), valueArgs...)
	if err != nil {
		return
	}
	medias, err = repository.scanRows(rows)
	if err != nil {
		return
	}
	return
}

// Insert inserts medias into the database
// and save the media file to the storage based on specified media disk (filesystem disk)
func (repository *MediaRepository) Insert(ctx context.Context, medias ...*entities.Media) (
	result sql.Result, err error) {
	var (
		// Query string
		query = buildBatchInsertQuery(repository.tableName, len(medias), repository.columnNames...)
		// Query value args
		valueArgs []any
		// saved media files paths
		savedFilePaths []string
		// Error Channel
		errChan = synch.MakeChanWithValue[error](nil, 1)
	)
	defer close(errChan)

	for _, media := range medias {
		if err != nil {
			fileh.BatchRemove(savedFilePaths...) // rollback saved files
			return
		}
		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			fileh.BatchRemove(savedFilePaths...) // rollback saved files
			return
		}

		repository.waitGroup.Add(1)

		go func(wg *sync.WaitGroup, media *entities.Media) {
			defer wg.Done()

			errChanVal := synch.GetChanValAndKeep(errChan)
			if errChanVal != nil {
				return
			}

			// check if file is nil or not provided if meet the condition return an error
			if media.File == nil || !media.File.IsProvided() {
				errChanVal = exceptions.FileNotProvided
				errChan <- errChanVal
				return
			}

			if media.Id == uuid.Nil {
				media.Id = uuid.New()
			}
			if media.CreatedAt.IsZero() {
				media.CreatedAt = time.Now()
			}
			if media.CollectionName == "" {
				media.CollectionName = "default"
			}
			if media.FileName == "" {
				media.FileName = filepath.Base(media.File.Name)
			}
			if media.Size == 0 {
				media.Size = media.File.Size
			}
			if media.MimeType == "" {
				media.MimeType = media.File.Mime.String()
			}
			if media.Disk == "" {
				media.Disk = media.GetDefaultDisk()
			}

			// Save media file
			errChanVal = media.SaveFile()
			if errChanVal != nil {
				errChan <- errChanVal
				return
			}

			repository.mutex.Lock()
			defer repository.mutex.Unlock()
			savedFilePaths = append(savedFilePaths, media.GetFilePath()) // assign media file path to "savedFilePaths" in case of error happen it will used for rollbacking the media saved files
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

		}(repository.waitGroup, media)
	}

	repository.waitGroup.Wait()

	if err != nil {
		return
	}
	err = <-errChan
	if err != nil {
		return
	}

	result, err = repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		fileh.BatchRemove(savedFilePaths...) // rollback saved files
		return
	}
	err = nil
	return
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

	// if file is nil or not provided its mean no file changes in this update
	if media.File == nil || !media.File.IsProvided() {
		err := media.RenameFile() // rename file incase of media.Filename is updated
		if err != nil {
			return nil, err
		}
	} else { // check if file is provided
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

// DeleteByIds deletes the entities associated with the given ids
func (repository *MediaRepository) DeleteByIds(ctx context.Context, ids ...string) (sql.Result, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query, valueArgs := buildBatchDeleteQueryByIds(repository.tableName, ids...)

	result, err := repository.db.ExecContext(ctx, query, valueArgs...)

	return result, err
}
