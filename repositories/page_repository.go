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

type PageRepository struct {
	db     *sql.DB
	Query  query_reqs.Query
	entity entities.Page
}

func NewPageRepository(db *sql.DB) *PageRepository {
	return &PageRepository{
		db:     db,
		Query:  query_reqs.Default(),
		entity: entities.Page{},
	}
}

/*
 * ----------------------------------------------------------------
 * Repository utilty methods ⬇
 * ----------------------------------------------------------------
 */

// scanRows scans rows of the database and returns it as structs, and returns error if any error has occurred.
func (repository *PageRepository) scanRows(rows *sql.Rows) (
	pages []entities.Page, err error) {
	for rows.Next() {
		page := entities.Page{}
		err := rows.Scan(
			&page.Id,
			&page.SpaceId,
			&page.Title,
			&page.Order,
			&page.CreatedAt,
			&page.UpdatedAt,
		)
		if err != nil {
			return pages, err
		}
		pages = append(pages, page)
	}

	return pages, nil
}

// scanRow scans only a row of the database and returns it as struct, and returns error if any error has occurred.
func (repository *PageRepository) scanRow(rows *sql.Rows) (page entities.Page, err error) {
	pages, err := repository.scanRows(rows)
	if err != nil {
		return
	}
	if len(pages) == 0 {
		return page, nil
	}
	page = pages[0]
	return page, nil
}

/*
 * ----------------------------------------------------------------
 * Repository query methods ⬇
 * ----------------------------------------------------------------
 */

func (repository *PageRepository) All(ctx context.Context) (pages []entities.Page, err error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.entity.GetColumnNames()) +
		" FROM " + repository.entity.GetTableName()
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()
	return repository.scanRows(rows)
}

// Find finds by id
func (repository *PageRepository) Find(ctx context.Context, id string) (pages entities.Page, err error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.entity.GetColumnNames()) +
		" FROM " + repository.entity.GetTableName() + " WHERE `id` = ? LIMIT 1"
	rows, err := repository.db.QueryContext(ctx, query, id)
	if err != nil {
		return
	}
	defer rows.Close()
	return repository.scanRow(rows)
}

// GetBySpaceId get pages by space id
func (repository *PageRepository) GetBySpaceId(ctx context.Context, spaceId string) (
	pages []entities.Page, err error) {

	var valueArgs []any
	queryBuf := bytes.NewBufferString("SELECT ")
	queryBuf.WriteString(stringh.SliceTableColumnToStr(
		repository.entity.GetTableName(),
		repository.entity.GetColumnNames(),
	))
	queryBuf.WriteString(" FROM ")
	queryBuf.WriteString(repository.entity.GetTableName())
	queryBuf.WriteString(" WHERE ")
	queryBuf.WriteString(repository.entity.GetTableName() + ".`space_id` = ?")
	valueArgs = append(valueArgs, spaceId)
	// TODO: implements pages search query feature
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

func (repository *PageRepository) Insert(ctx context.Context, pages ...*entities.Page) (sql.Result, error) {
	query := buildBatchInsertQuery(
		repository.entity.GetTableName(),
		len(pages),
		repository.entity.GetColumnNames()...,
	)
	valueArgs := []any{}
	for _, page := range pages {
		if page.Id == uuid.Nil {
			page.Id = uuid.New()
		}
		if page.CreatedAt.IsZero() {
			page.CreatedAt = time.Now()
		}
		valueArgs = append(valueArgs,
			page.Id,
			page.SpaceId,
			page.Title,
			page.Order,
			page.CreatedAt,
			page.UpdatedAt,
		)
	}
	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (repository *PageRepository) Create(ctx context.Context, page *entities.Page) (sql.Result, error) {
	return repository.Insert(ctx, page)
}
