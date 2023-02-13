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

type SpaceRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewSpaceRepo(db *sql.DB) *SpaceRepo {
	return &SpaceRepo{
		db:          db,
		tableName:   "spaces",
		columnNames: helper.GetStructFieldJsonTag(entities.Space{}),
	}
}

func (repo *SpaceRepo) All(ctx context.Context) ([]entities.Space, error) {
	query := "SELECT " + helper.DBSliceColumnsToStr(repo.columnNames) + " FROM " + repo.tableName
	spaces := []entities.Space{}
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
		return spaces, err
	}

	for rows.Next() {
		space := entities.Space{}
		err := rows.Scan(
			&space.Id, &space.Name, &space.Description,
			&space.Domain, &space.CreatedAt, &space.UpdatedAt,
		)
		if err != nil {
			helper.LogIfError(err)
			return spaces, err
		}
		spaces = append(spaces, space)
	}

	if len(spaces) == 0 {
		return spaces, exceptions.DataNotFoundError
	}

	return spaces, nil
}

func (repo *SpaceRepo) Insert(ctx context.Context, spaces ...entities.Space) ([]entities.Space, error) {
	query := buildBatchInsertQuery(repo.tableName, len(spaces), repo.columnNames...)
	valueArgs := []any{}

	for _, space := range spaces {
		if space.Id.String() == "" {
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

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
		return spaces, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.LogIfError(err)
		return spaces, err
	}
	return spaces, nil
}

func (repo *SpaceRepo) Create(ctx context.Context, space entities.Space) (entities.Space, error) {
	spaces, err := repo.Insert(ctx, space)
	if err != nil {
		return entities.Space{}, err
	}

	return spaces[0], nil
}
