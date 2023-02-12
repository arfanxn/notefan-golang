package repositories

import "database/sql"

type PageContentRepo struct {
	tableName string
	db        *sql.DB
}

func NewPageContentRepo(db *sql.DB) *PageContentRepo {
	return &PageContentRepo{
		tableName: "page_contents",
		db:        db,
	}
}
