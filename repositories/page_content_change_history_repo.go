package repositories

import "database/sql"

type PageContentChangeHistoryRepo struct {
	tableName string
	db        *sql.DB
}

func NewPageContentChangeHistoryRepo(db *sql.DB) *PageContentChangeHistoryRepo {
	return &PageContentChangeHistoryRepo{
		tableName: "page_content_change_history",
		db:        db,
	}
}
