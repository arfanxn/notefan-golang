package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"

	"github.com/google/uuid"
)

type PageContentRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewPageContentRepository(db *sql.DB) *PageContentRepository {
	return &PageContentRepository{
		db:          db,
		tableName:   "page_contents",
		columnNames: reflecth.GetFieldJsonTag(entities.PageContent{}),
	}
}

func (repository *PageContentRepository) All(ctx context.Context) ([]entities.PageContent, error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.columnNames) + " FROM " + repository.tableName
	pageContents := []entities.PageContent{}
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
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
			return pageContents, err
		}
		pageContents = append(pageContents, pageContent)
	}
	return pageContents, nil
}

func (repository *PageContentRepository) Insert(ctx context.Context, pageContents ...*entities.PageContent) (
	sql.Result, error) {
	query := buildBatchInsertQuery(repository.tableName, len(pageContents), repository.columnNames...)
	valueArgs := []any{}

	for _, pageContent := range pageContents {
		if pageContent.Id == uuid.Nil {
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
	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (repository *PageContentRepository) Create(ctx context.Context, pageContent *entities.PageContent) (
	sql.Result, error) {
	return repository.Insert(ctx, pageContent)
}
