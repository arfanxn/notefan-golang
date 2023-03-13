package repositories

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/query_reqs"
)

type PageContentChangeHistoryRepository struct {
	db     *sql.DB
	Query  query_reqs.Query
	entity entities.PageContentChangeHistory
}

func NewPageContentChangeHistoryRepository(db *sql.DB) *PageContentChangeHistoryRepository {
	return &PageContentChangeHistoryRepository{
		db:     db,
		Query:  query_reqs.Default(),
		entity: entities.PageContentChangeHistory{},
	}
}

func (repository *PageContentChangeHistoryRepository) All(ctx context.Context) ([]entities.PageContentChangeHistory, error) {
	query := "SELECT " + strings.Join(repository.entity.GetColumnNames(), ", ") +
		" FROM " + repository.entity.GetTableName()
	pageContentChangeHistories := []entities.PageContentChangeHistory{}
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		return pageContentChangeHistories, err
	}
	for rows.Next() {
		pcch := entities.PageContentChangeHistory{}
		err := rows.Scan(
			&pcch.BeforePageContentId,
			&pcch.AfterPageContentId,
			&pcch.UserId,
			&pcch.CreatedAt,
			&pcch.UpdatedAt,
		)
		if err != nil {
			return pageContentChangeHistories, err
		}
		pageContentChangeHistories = append(pageContentChangeHistories, pcch)
	}
	return pageContentChangeHistories, nil
}

func (repository *PageContentChangeHistoryRepository) Insert(
	ctx context.Context, spaces ...*entities.PageContentChangeHistory) (
	sql.Result, error) {
	query := buildBatchInsertQuery(
		repository.entity.GetTableName(),
		len(spaces),
		repository.entity.GetColumnNames()...,
	)
	valueArgs := []any{}
	for _, pcch := range spaces {
		if pcch.CreatedAt.IsZero() {
			pcch.CreatedAt = time.Now()
		}
		valueArgs = append(valueArgs,
			&pcch.BeforePageContentId,
			&pcch.AfterPageContentId,
			&pcch.UserId,
			&pcch.CreatedAt,
			&pcch.UpdatedAt,
		)
	}
	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (repository *PageContentChangeHistoryRepository) Create(ctx context.Context, pcch *entities.PageContentChangeHistory) (sql.Result, error) {
	return repository.Insert(ctx, pcch)
}
