package repositories

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/notefan-golang/config"
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

func (repository *MediaRepository) All(ctx context.Context) ([]entities.Media, error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.columnNames) + " FROM " + repository.tableName
	medias := []entities.Media{}
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		errorh.Log(err)
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
			errorh.Log(err)
			return medias, err
		}
		medias = append(medias, media)
	}

	if len(medias) == 0 {
		return medias, exceptions.HTTPNotFound
	}

	return medias, nil
}

// Insert inserts medias metadata into the database
// and save the media file to the storage based on specified media disk (filesystem disk)
func (repository *MediaRepository) Insert(ctx context.Context, medias ...entities.Media) ([]entities.Media, error) {
	query := buildBatchInsertQuery(repository.tableName, len(medias), repository.columnNames...)
	valueArgs := []any{}

	var err error
	var savedFilePaths []string

	for _, media := range medias {
		if err != nil {
			fileh.BatchRemove(savedFilePaths...) // rollback saved files
			return medias, err
		}

		repository.waitGroup.Add(1)

		go func(wg *sync.WaitGroup, media entities.Media) {
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

			media.GuessMimeType()

			// save file path destination
			path := filepath.Join(
				config.FSDisks[media.Disk].Root,
				"medias",
				media.Id.String(),
				filepath.Base(media.FileName),
			)

			// check if file exists, if not exists return an error
			if media.File.Len() <= 0 {
				err = exceptions.FileNotProvided
				return
			}

			// Save media file
			err = fileh.Save(path, media.File)

			repository.mutex.Lock()
			savedFilePaths = append(savedFilePaths, path) // assign saved file path to "savedFilePaths" in case of error happen it will used for rollbacking the saved files
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

	_, err = repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		errorh.Log(err)
		fileh.BatchRemove(savedFilePaths...) // rollback saved files
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
