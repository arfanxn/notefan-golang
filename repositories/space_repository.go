package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"

	"github.com/google/uuid"
)

type SpaceRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewSpaceRepository(db *sql.DB) *SpaceRepository {
	return &SpaceRepository{
		db:          db,
		tableName:   "spaces",
		columnNames: reflecth.GetFieldJsonTag(entities.Space{}),
	}
}

func (repository *SpaceRepository) All(ctx context.Context) ([]entities.Space, error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.columnNames) + " FROM " + repository.tableName
	spaces := []entities.Space{}
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		errorh.Log(err)
		return spaces, err
	}

	for rows.Next() {
		space := entities.Space{}
		err := rows.Scan(
			&space.Id, &space.Name, &space.Description,
			&space.Domain, &space.CreatedAt, &space.UpdatedAt,
		)
		if err != nil {
			errorh.Log(err)
			return spaces, err
		}
		spaces = append(spaces, space)
	}

	if len(spaces) == 0 {
		return spaces, exceptions.HTTPNotFound
	}

	return spaces, nil
}

func (repository *SpaceRepository) Insert(ctx context.Context, spaces ...entities.Space) ([]entities.Space, error) {
	query := buildBatchInsertQuery(repository.tableName, len(spaces), repository.columnNames...)
	valueArgs := []any{}

	for _, space := range spaces {
		if space.Id == uuid.Nil {
			space.Id = uuid.New()
		}
		if space.CreatedAt.IsZero() {
			space.CreatedAt = time.Now()
		}
		valueArgs = append(valueArgs,
			space.Id,
			space.Name,
			space.Description,
			space.Domain,
			space.CreatedAt,
			space.UpdatedAt,
		)
	}

	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		errorh.Log(err)
		return spaces, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		errorh.Log(err)
		return spaces, err
	}
	return spaces, nil
}

func (repository *SpaceRepository) Create(ctx context.Context, space entities.Space) (entities.Space, error) {
	spaces, err := repository.Insert(ctx, space)
	if err != nil {
		return entities.Space{}, err
	}

	return spaces[0], nil
}
