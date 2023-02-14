package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"strings"
	"time"
)

type PageContentChangeHistoryRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewPageContentChangeHistoryRepository(db *sql.DB) *PageContentChangeHistoryRepository {
	return &PageContentChangeHistoryRepository{
		db:          db,
		tableName:   "page_content_change_history",
		columnNames: helper.ReflectGetStructFieldJsonTag(entities.PageContentChangeHistory{}),
	}
}

func (repository *PageContentChangeHistoryRepository) All(ctx context.Context) ([]entities.PageContentChangeHistory, error) {
	query := "SELECT " + strings.Join(repository.columnNames, ", ") + " FROM " + repository.tableName
	pageContentChangeHistories := []entities.PageContentChangeHistory{}
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		helper.ErrorLog(err)
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
			helper.ErrorLog(err)
			return pageContentChangeHistories, err
		}
		pageContentChangeHistories = append(pageContentChangeHistories, pcch)
	}

	if len(pageContentChangeHistories) == 0 {
		return pageContentChangeHistories, exceptions.HTTPNotFound
	}

	return pageContentChangeHistories, nil
}

func (repository *PageContentChangeHistoryRepository) Insert(
	ctx context.Context, spaces ...entities.PageContentChangeHistory) (
	[]entities.PageContentChangeHistory, error) {
	query := buildBatchInsertQuery(repository.tableName, len(spaces), repository.columnNames...)
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

	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		helper.ErrorLog(err)
		return spaces, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.ErrorLog(err)
		return spaces, err
	}
	return spaces, nil
}

func (repository *PageContentChangeHistoryRepository) Create(ctx context.Context, space entities.PageContentChangeHistory) (entities.PageContentChangeHistory, error) {
	spaces, err := repository.Insert(ctx, space)
	if err != nil {
		return entities.PageContentChangeHistory{}, err
	}

	return spaces[0], nil
}
