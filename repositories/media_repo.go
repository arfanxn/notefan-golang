package repositories

import "database/sql"

type MediaRepo struct {
	tableName string
	db        *sql.DB
}

func NewMediaRepo(db *sql.DB) *MediaRepo {
	return &MediaRepo{
		tableName: "medias",
		db:        db,
	}
}
