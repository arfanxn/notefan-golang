package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"

	"github.com/google/uuid"
)

type PageContentRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewPageContentRepo(db *sql.DB) *PageContentRepo {
	return &PageContentRepo{
		db:          db,
		tableName:   "page_contents",
		columnNames: helper.GetStructFieldJsonTag(entities.PageContent{}),
	}
}

func (repo *PageContentRepo) All(ctx context.Context) ([]entities.PageContent, error) {
	query := "SELECT " + helper.DBSliceColumnsToStr(repo.columnNames) + " FROM " + repo.tableName
	pageContents := []entities.PageContent{}
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
		return pageContents, err
	}

	for rows.Next() {
		pageContent := entities.PageContent{}
		err := rows.Scan(
			&pageContent.Id,
			&pageContent.PageId,
			&pageContent.Type,
			&pageContent.Order,
			&pageContent.Body,
			&pageContent.CreatedAt,
			&pageContent.UpdatedAt,
		)
		if err != nil {
			helper.LogIfError(err)
			return pageContents, err
		}
		pageContents = append(pageContents, pageContent)
	}

	if len(pageContents) == 0 {
		return pageContents, exceptions.DataNotFoundError
	}

	return pageContents, nil
}

func (repo *PageContentRepo) Insert(ctx context.Context, pageContents ...entities.PageContent) ([]entities.PageContent, error) {
	query := buildBatchInsertQuery(repo.tableName, len(pageContents), repo.columnNames...)
	valueArgs := []any{}

	for _, pageContent := range pageContents {
		if pageContent.Id.String() == "" {
			pageContent.Id = uuid.New()
		}
		if pageContent.CreatedAt.IsZero() {
			pageContent.CreatedAt = time.Now()
		}
		valueArgs = append(valueArgs,
			pageContent.Id,
			pageContent.PageId,
			pageContent.Type,
			pageContent.Order,
			pageContent.Body,
			pageContent.CreatedAt,
			pageContent.UpdatedAt,
		)
	}

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
		return pageContents, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.LogIfError(err)
		return pageContents, err
	}
	return pageContents, nil
}

func (repo *PageContentRepo) Create(ctx context.Context, pageContent entities.PageContent) (entities.PageContent, error) {
	pageContents, err := repo.Insert(ctx, pageContent)
	if err != nil {
		return entities.PageContent{}, err
	}

	return pageContents[0], nil
}
