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

type PageContentChangeHistoryRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewPageContentChangeHistoryRepo(db *sql.DB) *PageContentChangeHistoryRepo {
	return &PageContentChangeHistoryRepo{
		db:          db,
		tableName:   "page_content_change_history",
		columnNames: helper.ReflectGetStructFieldJsonTag(entities.PageContentChangeHistory{}),
	}
}

func (repo *PageContentChangeHistoryRepo) All(ctx context.Context) ([]entities.PageContentChangeHistory, error) {
	query := "SELECT " + strings.Join(repo.columnNames, ", ") + " FROM " + repo.tableName
	pageContentChangeHistories := []entities.PageContentChangeHistory{}
	rows, err := repo.db.QueryContext(ctx, query)
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
		return pageContentChangeHistories, exceptions.DataNotFoundError
	}

	return pageContentChangeHistories, nil
}

func (repo *PageContentChangeHistoryRepo) Insert(
	ctx context.Context, spaces ...entities.PageContentChangeHistory) (
	[]entities.PageContentChangeHistory, error) {
	query := buildBatchInsertQuery(repo.tableName, len(spaces), repo.columnNames...)
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

	stmt, err := repo.db.PrepareContext(ctx, query)
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

func (repo *PageContentChangeHistoryRepo) Create(ctx context.Context, space entities.PageContentChangeHistory) (entities.PageContentChangeHistory, error) {
	spaces, err := repo.Insert(ctx, space)
	if err != nil {
		return entities.PageContentChangeHistory{}, err
	}

	return spaces[0], nil
}
