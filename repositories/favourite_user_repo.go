package repositories

import "database/sql"

type FavouriteUserRepo struct {
	tableName string
	db        *sql.DB
}

func NewFavouriteUserRepo(db *sql.DB) *FavouriteUserRepo {
	return &FavouriteUserRepo{
		tableName: "favourite_user",
		db:        db,
	}
}
