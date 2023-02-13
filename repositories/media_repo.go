package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"

	"github.com/google/uuid"
)

type MediaRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewMediaRepo(db *sql.DB) *MediaRepo {
	return &MediaRepo{
		db:          db,
		tableName:   "medias",
		columnNames: helper.ReflectGetStructFieldJsonTag(entities.Media{}),
	}
}

func (repo *MediaRepo) All(ctx context.Context) ([]entities.Media, error) {
	query := "SELECT " + helper.SliceTableColumnsToString(repo.columnNames) + " FROM " + repo.tableName
	medias := []entities.Media{}
	rows, err := repo.db.QueryContext(ctx, query)
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
		return medias, exceptions.DataNotFoundError
	}

	return medias, nil
}

func (repo *MediaRepo) Insert(ctx context.Context, medias ...entities.Media) ([]entities.Media, error) {
	query := buildBatchInsertQuery(repo.tableName, len(medias), repo.columnNames...)
	valueArgs := []any{}

	for _, media := range medias {
		if media.Id.String() == "" {
			media.Id = uuid.New()
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

	stmt, err := repo.db.PrepareContext(ctx, query)
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

func (repo *MediaRepo) Create(ctx context.Context, media entities.Media) (entities.Media, error) {
	medias, err := repo.Insert(ctx, media)
	if err != nil {
		return entities.Media{}, err
	}

	return medias[0], nil
}
