package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"
)

type FavouriteUserRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewFavouriteUserRepo(db *sql.DB) *FavouriteUserRepo {
	return &FavouriteUserRepo{
		db:          db,
		tableName:   "favourite_user",
		columnNames: helper.GetStructFieldJsonTag(entities.FavouriteUser{}),
	}
}

func (repo *FavouriteUserRepo) All(ctx context.Context) ([]entities.FavouriteUser, error) {
	query := "SELECT " + helper.DBSliceColumnsToStr(repo.columnNames) + " FROM " + repo.tableName
	favouriteUsers := []entities.FavouriteUser{}
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
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
			helper.LogIfError(err)
			return favouriteUsers, err
		}
		favouriteUsers = append(favouriteUsers, favouriteUser)
	}

	if len(favouriteUsers) == 0 {
		return favouriteUsers, exceptions.DataNotFoundError
	}

	return favouriteUsers, nil
}

func (repo *FavouriteUserRepo) Insert(ctx context.Context, favouriteUsers ...entities.FavouriteUser) ([]entities.FavouriteUser, error) {
	query := buildBatchInsertQuery(repo.tableName, len(favouriteUsers), repo.columnNames...)
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

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
		return favouriteUsers, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.LogIfError(err)
		return favouriteUsers, err
	}
	return favouriteUsers, nil
}

func (repo *FavouriteUserRepo) Create(ctx context.Context, favouriteUser entities.FavouriteUser) (entities.FavouriteUser, error) {
	favouriteUsers, err := repo.Insert(ctx, favouriteUser)
	if err != nil {
		return entities.FavouriteUser{}, err
	}

	return favouriteUsers[0], nil
}
