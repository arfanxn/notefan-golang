package repositories

import (
	"bytes"
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

/*
 * ----------------------------------------------------------------
 * Repository utilty methods ⬇
 * ----------------------------------------------------------------
 */

// scanRows scans rows of the database and returns it as structs, and returns error if any error has occurred.
func (repository *SpaceRepository) scanRows(rows *sql.Rows) (spaces []entities.Space, err error) {
	for rows.Next() {
		space := entities.Space{}
		err := rows.Scan(
			&space.Id,
			&space.Name,
			&space.Description,
			&space.Domain,
			&space.CreatedAt,
			&space.UpdatedAt,
		)
		errorh.LogPanic(err) // panic if scan fails
		spaces = append(spaces, space)
	}

	if len(spaces) == 0 {
		return spaces, exceptions.HTTPNotFound
	}
	return spaces, nil
}

// scanRow scans only a row of the database and returns it as struct, and returns error if any error has occurred.
func (repository *SpaceRepository) scanRow(rows *sql.Rows) (entities.Space, error) {
	spaces, err := repository.scanRows(rows)
	if err != nil {
		return entities.Space{}, err
	}
	return spaces[0], nil
}

/*
 * ----------------------------------------------------------------
 * Repository query methods ⬇
 * ----------------------------------------------------------------
 */

// All retrieves all data on table from database
func (repository *SpaceRepository) All(ctx context.Context) ([]entities.Space, error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.columnNames) + " FROM " + repository.tableName
	rows, err := repository.db.QueryContext(ctx, query)
	errorh.LogPanic(err)
	return repository.scanRows(rows)
}

// Find retrieves data on table from database by the given id
func (repository *SpaceRepository) Find(ctx context.Context, id string) (entities.Space, error) {
	queryBuf := bytes.NewBufferString("SELECT ")
	queryBuf.WriteString(stringh.SliceColumnToStr(repository.columnNames))
	queryBuf.WriteString(" FROM ")
	queryBuf.WriteString(repository.tableName)
	queryBuf.WriteString(" WHERE `id` = ?")
	rows, err := repository.db.QueryContext(ctx, queryBuf.String(), id)
	errorh.LogPanic(err)
	return repository.scanRow(rows)
}

// Insert inserts into database
func (repository *SpaceRepository) Insert(ctx context.Context, spaces ...*entities.Space) (sql.Result, error) {
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

	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	errorh.Log(err)
	return result, err
}

// Create creates and save into database
func (repository *SpaceRepository) Create(ctx context.Context, space *entities.Space) (sql.Result, error) {
	result, err := repository.Insert(ctx, space)
	return result, err
}

// UpdateById updates entity by id
func (repository *SpaceRepository) UpdateById(ctx context.Context, space *entities.Space) (sql.Result, error) {
	query := buildUpdateQuery(repository.tableName, repository.columnNames...) + " WHERE id = ?"

	// Refresh entity updated at
	space.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	result, err := repository.db.ExecContext(ctx, query,
		space.Id,
		space.Name,
		space.Description,
		space.Domain,
		space.CreatedAt,
		space.UpdatedAt,
		space.Id)

	return result, err
}

// DeleteByIds deletes entities by the given ids
func (repository *SpaceRepository) DeleteByIds(ctx context.Context, ids ...string) (sql.Result, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query, valueArgs := buildBatchDeleteQueryByIds(repository.tableName, ids...)

	result, err := repository.db.ExecContext(ctx, query, valueArgs...)

	return result, err
}
