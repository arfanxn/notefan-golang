package repositories

import (
	"bytes"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/query_reqs"
)

type UserRoleSpaceRepository struct {
	db     *sql.DB
	Query  query_reqs.Query
	entity entities.UserRoleSpace
}

func NewUserRoleSpaceRepository(db *sql.DB) *UserRoleSpaceRepository {
	return &UserRoleSpaceRepository{
		db:     db,
		Query:  query_reqs.Default(),
		entity: entities.UserRoleSpace{},
	}
}

/*
 * ----------------------------------------------------------------
 * Repository utilty methods ⬇
 * ----------------------------------------------------------------
 */

// scanRows scans rows of the database and returns it as structs, and returns error if any error has occurred.
func (repository *UserRoleSpaceRepository) scanRows(rows *sql.Rows) (
	urss []entities.UserRoleSpace, err error) {
	for rows.Next() {
		urs := entities.UserRoleSpace{}
		err := rows.Scan(
			&urs.Id,
			&urs.UserId,
			&urs.RoleId,
			&urs.SpaceId,
			&urs.CreatedAt,
			&urs.UpdatedAt,
		)
		if err != nil {
			return urss, err
		}
		urss = append(urss, urs)
	}

	return urss, nil
}

// scanRow scans only a row of the database and returns it as struct, and returns error if any error has occurred.
func (repository *UserRoleSpaceRepository) scanRow(rows *sql.Rows) (urs entities.UserRoleSpace, err error) {
	urss, err := repository.scanRows(rows)
	if err != nil {
		return
	}
	if len(urss) == 0 {
		return urs, nil
	}
	urs, err = urss[0], nil
	return
}

/*
 * ----------------------------------------------------------------
 * Repository query methods ⬇
 * ----------------------------------------------------------------
 */

// FindByUserIdAndSpaceId finds row by user id and space id
func (repository *UserRoleSpaceRepository) FindByUserIdAndSpaceId(
	ctx context.Context, userId string, spaceId string) (
	urs entities.UserRoleSpace, err error) {
	queryBuf := bytes.NewBufferString("SELECT ")
	queryBuf.WriteString(stringh.SliceColumnToStr(repository.entity.GetColumnNames()))
	queryBuf.WriteString(" FROM ")
	queryBuf.WriteString(repository.entity.GetTableName())
	queryBuf.WriteString(" WHERE `user_id` = ? AND `space_id` = ?")
	rows, err := repository.db.QueryContext(ctx, queryBuf.String(), userId, spaceId)
	if err != nil {
		return
	}
	urs, err = repository.scanRow(rows)
	if err != nil {
		return
	}
	return
}

func (repository *UserRoleSpaceRepository) Insert(ctx context.Context, userRoleSpaces ...*entities.UserRoleSpace) (
	sql.Result, error) {
	query := buildBatchInsertQuery(
		repository.entity.GetTableName(),
		len(userRoleSpaces),
		repository.entity.GetColumnNames()...,
	)
	valueArgs := []any{}

	for _, userRoleSpace := range userRoleSpaces {
		if userRoleSpace.Id == uuid.Nil {
			userRoleSpace.Id = uuid.New()
		}
		if userRoleSpace.CreatedAt.IsZero() {
			userRoleSpace.CreatedAt = time.Now()
		}
		valueArgs = append(valueArgs,
			userRoleSpace.Id,
			userRoleSpace.UserId,
			userRoleSpace.RoleId,
			userRoleSpace.SpaceId,
			userRoleSpace.CreatedAt,
			userRoleSpace.UpdatedAt,
		)
	}

	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	return result, err
}

func (repository *UserRoleSpaceRepository) Create(ctx context.Context, userRoleSpace *entities.UserRoleSpace) (
	sql.Result, error) {
	result, err := repository.Insert(ctx, userRoleSpace)
	if err != nil {
		return result, err
	}

	return result, nil
}

// DeleteByIds deletes the entities associated with the given ids
func (repository *UserRoleSpaceRepository) DeleteByIds(ctx context.Context, ids ...string) (sql.Result, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query, valueArgs := buildBatchDeleteQueryByIds(repository.entity.GetTableName(), ids...)
	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	return result, err
}
