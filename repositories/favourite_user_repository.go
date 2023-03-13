package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/query_reqs"
)

type FavouriteUserRepository struct {
	db     *sql.DB
	Query  query_reqs.Query
	entity entities.FavouriteUser
}

func NewFavouriteUserRepository(db *sql.DB) *FavouriteUserRepository {
	return &FavouriteUserRepository{
		db:     db,
		Query:  query_reqs.Default(),
		entity: entities.FavouriteUser{},
	}
}

func (repository *FavouriteUserRepository) All(ctx context.Context) ([]entities.FavouriteUser, error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.entity.GetColumnNames()) +
		" FROM " + repository.entity.GetTableName()
	favouriteUsers := []entities.FavouriteUser{}
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		return favouriteUsers, err
	}
	for rows.Next() {
		favouriteUser := entities.FavouriteUser{}
		err := rows.Scan(
			&favouriteUser.FavouriteableType,
			&favouriteUser.FavouriteableId,
			&favouriteUser.UserId,
			&favouriteUser.Order,
			&favouriteUser.CreatedAt,
			&favouriteUser.UpdatedAt,
		)
		if err != nil {
			return favouriteUsers, err
		}
		favouriteUsers = append(favouriteUsers, favouriteUser)
	}
	return favouriteUsers, nil
}

func (repository *FavouriteUserRepository) Insert(ctx context.Context, favouriteUsers ...*entities.FavouriteUser) (sql.Result, error) {
	query := buildBatchInsertQuery(
		repository.entity.GetTableName(),
		len(favouriteUsers),
		repository.entity.GetColumnNames()...,
	)
	valueArgs := []any{}
	for _, favouriteUser := range favouriteUsers {
		if favouriteUser.CreatedAt.IsZero() {
			favouriteUser.CreatedAt = time.Now()
		}
		valueArgs = append(valueArgs,
			favouriteUser.FavouriteableType,
			favouriteUser.FavouriteableId,
			favouriteUser.UserId,
			favouriteUser.Order,
			favouriteUser.CreatedAt,
			favouriteUser.UpdatedAt,
		)
	}
	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (repository *FavouriteUserRepository) Create(ctx context.Context, favouriteUser *entities.FavouriteUser) (
	sql.Result, error) {
	return repository.Insert(ctx, favouriteUser)
}
