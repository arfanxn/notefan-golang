package repositories

import (
	"bytes"
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/query_reqs"

	"github.com/google/uuid"
)

type SpaceRepository struct {
	db     *sql.DB
	Query  query_reqs.Query
	entity entities.Space
}

func NewSpaceRepository(db *sql.DB) *SpaceRepository {
	return &SpaceRepository{
		db:     db,
		Query:  query_reqs.Default(),
		entity: entities.Space{},
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
		if err != nil {
			return spaces, err
		}
		spaces = append(spaces, space)
	}
	return spaces, nil
}

// scanRow scans only a row of the database and returns it as struct, and returns error if any error has occurred.
func (repository *SpaceRepository) scanRow(rows *sql.Rows) (space entities.Space, err error) {
	spaces, err := repository.scanRows(rows)
	if len(spaces) == 0 {
		return
	}
	space, err = spaces[0], nil
	return
}

/*
 * ----------------------------------------------------------------
 * Repository query methods ⬇
 * ----------------------------------------------------------------
 */

// All retrieves all data on table from database
func (repository *SpaceRepository) All(ctx context.Context) (spaces []entities.Space, err error) {
	query := "SELECT " +
		stringh.SliceColumnToStr(repository.entity.GetColumnNames()) + " FROM " + repository.entity.GetTableName()
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		return
	}
	spaces, err = repository.scanRows(rows)
	return
}

// Find retrieves data on table from database by the given id
func (repository *SpaceRepository) Find(ctx context.Context, id string) (space entities.Space, err error) {
	queryBuf := bytes.NewBufferString("SELECT ")
	queryBuf.WriteString(stringh.SliceColumnToStr(repository.entity.GetColumnNames()))
	queryBuf.WriteString(" FROM ")
	queryBuf.WriteString(repository.entity.GetTableName())
	queryBuf.WriteString(" WHERE `id` = ?")
	rows, err := repository.db.QueryContext(ctx, queryBuf.String(), id)
	if err != nil {
		return
	}
	space, err = repository.scanRow(rows)
	if err != nil {
		return
	}
	return space, err
}

// GetByUserId get spaces by user id
func (repository *SpaceRepository) GetByUserId(ctx context.Context, userId string) (
	spaces []entities.Space, err error) {

	var (
		ety    entities.Space
		ursEty entities.UserRoleSpace
	)

	var valueArgs []any
	queryBuf := bytes.NewBufferString("SELECT ")
	queryBuf.WriteString(stringh.SliceTableColumnToStr(
		ety.GetTableName(), ety.GetColumnNames())) // column names
	queryBuf.WriteRune(',') // comma
	queryBuf.WriteString(stringh.SliceTableColumnToStr(
		ursEty.GetTableName(), ursEty.GetColumnNames(), // join column names
	))
	queryBuf.WriteString(" FROM ")
	queryBuf.WriteString(repository.entity.GetTableName())
	queryBuf.WriteString(" INNER JOIN ")
	queryBuf.WriteString(ursEty.GetTableName())
	queryBuf.WriteString(" ON ")
	queryBuf.WriteString(repository.entity.GetTableName() + ".id")
	queryBuf.WriteString(" = ")
	queryBuf.WriteString(ursEty.GetTableName() + ".space_id")
	queryBuf.WriteString(" WHERE ")
	queryBuf.WriteString(ursEty.GetTableName() + ".user_id = ?")
	valueArgs = append(valueArgs, userId)
	// TODO: fix error on search by keyword
	if repository.Query.Keyword != "" { // if keyword is set then add to query
		queryBuf.WriteString(" AND ( ")
		queryBuf.WriteString(repository.entity.GetTableName() + ".name LIKE ?")
		queryBuf.WriteString(" OR ")
		queryBuf.WriteString(repository.entity.GetTableName() + ".description LIKE ?")
		queryBuf.WriteString(" OR ")
		queryBuf.WriteString(repository.entity.GetTableName() + ".domain LIKE ?")
		queryBuf.WriteString(" )")
		keyword := repository.Query.Keyword
		valueArgs = append(valueArgs, keyword, keyword, keyword)
	}
	queryBuf.WriteString(" LIMIT ? OFFSET ? ")
	valueArgs = append(valueArgs, repository.Query.Limit, repository.Query.Offset)
	if err != nil {
		return
	}

	rows, err := repository.db.QueryContext(ctx, queryBuf.String(), valueArgs...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var space entities.Space
		var urs entities.UserRoleSpace
		err := rows.Scan(
			&space.Id,
			&space.Name,
			&space.Description,
			&space.Domain,
			&space.CreatedAt,
			&space.UpdatedAt,
			&urs.Id,
			&urs.UserId,
			&urs.RoleId,
			&urs.SpaceId,
			&urs.CreatedAt,
			&urs.UpdatedAt,
		)
		if err != nil {
			return spaces, err
		}
		spaces = append(spaces, space)
	}
	return spaces, nil
}

// Insert inserts into database
func (repository *SpaceRepository) Insert(ctx context.Context, spaces ...*entities.Space) (sql.Result, error) {
	query := buildBatchInsertQuery(
		repository.entity.GetTableName(),
		len(spaces),
		repository.entity.GetColumnNames()...)
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
	if err != nil {
		return result, err
	}
	return result, err
}

// Create creates and save into database
func (repository *SpaceRepository) Create(ctx context.Context, space *entities.Space) (sql.Result, error) {
	return repository.Insert(ctx, space)
}

// UpdateById updates entity by id
func (repository *SpaceRepository) UpdateById(ctx context.Context, space *entities.Space) (sql.Result, error) {
	query := buildUpdateQuery(repository.entity.GetTableName(),
		repository.entity.GetColumnNames()...) + " WHERE id = ?"

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
	query, valueArgs := buildBatchDeleteQueryByIds(repository.entity.GetTableName(), ids...)
	return repository.db.ExecContext(ctx, query, valueArgs...)
}
