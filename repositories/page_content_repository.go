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
		errorh.Log(err)
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
			errorh.Log(err)
			return pageContents, err
		}
		pageContents = append(pageContents, pageContent)
	}

	if len(pageContents) == 0 {
		return pageContents, exceptions.HTTPNotFound
	}

	return pageContents, nil
}

func (repository *PageContentRepository) Insert(ctx context.Context, pageContents ...entities.PageContent) ([]entities.PageContent, error) {
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

	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		errorh.Log(err)
		return pageContents, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		errorh.Log(err)
		return pageContents, err
	}
	return pageContents, nil
}

func (repository *PageContentRepository) Create(ctx context.Context, pageContent entities.PageContent) (entities.PageContent, error) {
	pageContents, err := repository.Insert(ctx, pageContent)
	if err != nil {
		return entities.PageContent{}, err
	}

	return pageContents[0], nil
}
