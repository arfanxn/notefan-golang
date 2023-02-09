package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
)

type PageRepo struct {
	TableName string
	DBTX      *sql.DB
}

func NewPageRepo(dbtx *sql.DB) *PageRepo {
	return &PageRepo{
		TableName: "pages",
		DBTX:      dbtx,
	}
}

func (repo *PageRepo) scanRows(rows *sql.Rows) []entities.Page {
	var pages []entities.Page
	for rows.Next() {
		page := entities.Page{}
		err := rows.Scan(&page.Id, &page.Title)
		helper.LogFatalIfError(err)
		pages = append(pages, page)
	}
	return pages
}

func (repo *PageRepo) Get(ctx context.Context) []entities.Page {
	query := "SELECT id, title FROM " + repo.TableName
	rows, err := repo.DBTX.QueryContext(ctx, query)
	helper.LogFatalIfError(err)
	defer rows.Close()
	return repo.scanRows(rows)
}
