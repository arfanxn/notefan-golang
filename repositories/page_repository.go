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

func (repository *PageRepository) All(ctx context.Context) ([]entities.Page, error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.columnNames) + " FROM " + repository.tableName
	rows, err := repository.db.QueryContext(ctx, query)
	errorh.LogFatal(err)
	defer rows.Close()

	var pages []entities.Page
	for rows.Next() {
		page := entities.Page{}
		err := rows.Scan(&page.Id, &page.SpaceId, &page.Title, &page.Order, &page.CreatedAt, &page.UpdatedAt)
		errorh.LogFatal(err)
		pages = append(pages, page)
	}

	if len(pages) == 0 {
		return pages, exceptions.HTTPNotFound
	}

	return pages, nil
}

func (repository *PageRepository) Insert(ctx context.Context, pages ...entities.Page) ([]entities.Page, error) {
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

	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		errorh.Log(err)
		return pages, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		errorh.Log(err)
		return pages, err
	}
	return pages, nil
}

func (repository *PageRepository) Create(ctx context.Context, page entities.Page) (entities.Page, error) {
	pages, err := repository.Insert(ctx, page)
	if err != nil {
		return entities.Page{}, err
	}

	return pages[0], nil
}
