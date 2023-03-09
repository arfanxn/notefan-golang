package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"

	"github.com/google/uuid"
)

type PageRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewPageRepository(db *sql.DB) *PageRepository {
	return &PageRepository{
		db:          db,
		tableName:   "pages",
		columnNames: reflecth.GetFieldJsonTag(entities.Page{}),
	}
}

func (repository *PageRepository) All(ctx context.Context) (pages []entities.Page, err error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.columnNames) + " FROM " + repository.tableName
	rows, err := repository.db.QueryContext(ctx, query)
	errorh.Log(err)
	defer rows.Close()
	for rows.Next() {
		page := entities.Page{}
		err := rows.Scan(&page.Id, &page.SpaceId, &page.Title, &page.Order, &page.CreatedAt, &page.UpdatedAt)
		errorh.Log(err)
		pages = append(pages, page)
	}
	return pages, nil
}

func (repository *PageRepository) Insert(ctx context.Context, pages ...*entities.Page) (sql.Result, error) {
	query := buildBatchInsertQuery(repository.tableName, len(pages), repository.columnNames...)
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
		errorh.Log(err)
		return result, err
	}
	return result, nil
}

func (repository *PageRepository) Create(ctx context.Context, page *entities.Page) (sql.Result, error) {
	return repository.Insert(ctx, page)
}
