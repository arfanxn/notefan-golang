package repositories

import (
	"bytes"
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/query_reqs"

	"github.com/google/uuid"
)

type PageContentRepository struct {
	db     *sql.DB
	Query  query_reqs.Query
	entity entities.PageContent
}

func NewPageContentRepository(db *sql.DB) *PageContentRepository {
	return &PageContentRepository{
		db:     db,
		Query:  query_reqs.Default(),
		entity: entities.PageContent{},
	}
}

/*
 * ----------------------------------------------------------------
 * Repository utilty methods ⬇
 * ----------------------------------------------------------------
 */

// scanRows scans rows of the database and returns it as structs, and returns error if any error has occurred.
func (repository *PageContentRepository) scanRows(rows *sql.Rows) (
	pcs []entities.PageContent, err error) {
	for rows.Next() {
		pc := entities.PageContent{}
		err := rows.Scan(
			&pc.Id,
			&pc.PageId,
			&pc.Type,
			&pc.Order,
			&pc.Body,
			&pc.CreatedAt,
			&pc.UpdatedAt,
		)
		if err != nil {
			return pcs, err
		}
		pcs = append(pcs, pc)
	}
	return pcs, nil
}

// scanRow scans only a row of the database and returns it as struct, and returns error if any error has occurred.
func (repository *PageContentRepository) scanRow(rows *sql.Rows) (pc entities.PageContent, err error) {
	pcs, err := repository.scanRows(rows)
	if err != nil {
		return
	}
	if len(pcs) == 0 {
		return pc, nil
	}
	pc = pcs[0]
	return pc, nil
}

/*
 * ----------------------------------------------------------------
 * Repository query methods ⬇
 * ----------------------------------------------------------------
 */

func (repository *PageContentRepository) All(ctx context.Context) ([]entities.PageContent, error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.entity.GetColumnNames()) +
		" FROM " + repository.entity.GetTableName()
	pageContents := []entities.PageContent{}
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		return pageContents, err
	}
	return repository.scanRows(rows)
}

// Find finds by page id
func (repository *PageContentRepository) Find(ctx context.Context, id string) (
	pc entities.PageContent, err error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.entity.GetColumnNames()) +
		" FROM " + repository.entity.GetTableName() + " WHERE `id` = ? LIMIT 1"
	rows, err := repository.db.QueryContext(ctx, query, id)
	if err != nil {
		return
	}
	defer rows.Close()
	return repository.scanRow(rows)
}

// GetByPageId get pages by space id
func (repository *PageContentRepository) GetByPageId(ctx context.Context, pageId string) (
	pcs []entities.PageContent, err error) {
	var valueArgs []any
	queryBuf := bytes.NewBufferString("SELECT ")
	queryBuf.WriteString(stringh.SliceTableColumnToStr(
		repository.entity.GetTableName(),
		repository.entity.GetColumnNames(),
	))
	queryBuf.WriteString(" FROM ")
	queryBuf.WriteString(repository.entity.GetTableName())
	queryBuf.WriteString(" WHERE ")
	queryBuf.WriteString(repository.entity.GetTableName() + ".`page_id` = ?")
	valueArgs = append(valueArgs, pageId)
	queryBuf.WriteString(" LIMIT ? OFFSET ? ")
	valueArgs = append(valueArgs, repository.Query.Limit, repository.Query.Offset)
	if err != nil {
		return
	}
	rows, err := repository.db.QueryContext(ctx, queryBuf.String(), valueArgs...)
	if err != nil {
		return
	}
	defer rows.Close()
	return repository.scanRows(rows)
}

func (repository *PageContentRepository) Insert(ctx context.Context, pageContents ...*entities.PageContent) (
	sql.Result, error) {
	query := buildBatchInsertQuery(
		repository.entity.GetTableName(),
		len(pageContents),
		repository.entity.GetColumnNames()...,
	)
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

// UpdateById updates entity by id
func (repository *PageContentRepository) UpdateById(ctx context.Context, pageContent *entities.PageContent) (sql.Result, error) {
	query := buildUpdateQuery(repository.entity.GetTableName(),
		repository.entity.GetColumnNames()...) + " WHERE `id` = ?"

	// Refresh entity updated at
	pageContent.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	result, err := repository.db.ExecContext(ctx, query,
		pageContent.Id,
		pageContent.PageId,
		pageContent.Type,
		pageContent.Order,
		pageContent.Body,
		pageContent.CreatedAt,
		pageContent.UpdatedAt,
		pageContent.Id)

	return result, err
}

// DeleteByIds deletes entities by the given ids
func (repository *PageContentRepository) DeleteByIds(ctx context.Context, ids ...string) (sql.Result, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query, valueArgs := buildBatchDeleteQueryByIds(repository.entity.GetTableName(), ids...)
	return repository.db.ExecContext(ctx, query, valueArgs...)
}
