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
)

type FavouriteUserRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewFavouriteUserRepository(db *sql.DB) *FavouriteUserRepository {
	return &FavouriteUserRepository{
		db:          db,
		tableName:   "favourite_user",
		columnNames: reflecth.GetFieldJsonTag(entities.FavouriteUser{}),
	}
}

func (repository *FavouriteUserRepository) All(ctx context.Context) ([]entities.FavouriteUser, error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.columnNames) + " FROM " + repository.tableName
	favouriteUsers := []entities.FavouriteUser{}
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		errorh.Log(err)
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
			errorh.Log(err)
			return favouriteUsers, err
		}
		favouriteUsers = append(favouriteUsers, favouriteUser)
	}

	if len(favouriteUsers) == 0 {
		return favouriteUsers, exceptions.HTTPNotFound
	}

	return favouriteUsers, nil
}

func (repository *FavouriteUserRepository) Insert(ctx context.Context, favouriteUsers ...entities.FavouriteUser) ([]entities.FavouriteUser, error) {
	query := buildBatchInsertQuery(repository.tableName, len(favouriteUsers), repository.columnNames...)
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

	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		errorh.Log(err)
		return favouriteUsers, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		errorh.Log(err)
		return favouriteUsers, err
	}
	return favouriteUsers, nil
}

func (repository *FavouriteUserRepository) Create(ctx context.Context, favouriteUser entities.FavouriteUser) (entities.FavouriteUser, error) {
	favouriteUsers, err := repository.Insert(ctx, favouriteUser)
	if err != nil {
		return entities.FavouriteUser{}, err
	}

	return favouriteUsers[0], nil
}
