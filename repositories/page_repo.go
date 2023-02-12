package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"

	"github.com/google/uuid"
)

type PageRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewPageRepo(db *sql.DB) *PageRepo {
	return &PageRepo{
		db:          db,
		tableName:   "pages",
		columnNames: helper.GetStructFieldJsonTag(entities.Page{}),
	}
}

func (repo *PageRepo) Get(ctx context.Context) []entities.Page {
	// TODO: implement dynamic column names

	query := "SELECT id, title FROM " + repo.tableName
	rows, err := repo.db.QueryContext(ctx, query)
	helper.LogFatalIfError(err)
	defer rows.Close()

	var pages []entities.Page
	for rows.Next() {
		page := entities.Page{}
		err := rows.Scan(&page.Id, &page.Title)
		helper.LogFatalIfError(err)
		pages = append(pages, page)
	}
	return pages
}

func (repo *PageRepo) Insert(ctx context.Context, pages ...entities.Page) ([]entities.Page, error) {
	query := buildBatchInsertQuery(repo.tableName, len(pages), repo.columnNames...)
	valueArgs := []any{}

	for _, page := range pages {
		if page.Id.String() == "" {
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

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
		return pages, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.LogIfError(err)
		return pages, err
	}
	return pages, nil
}

func (repo *PageRepo) Create(ctx context.Context, page entities.Page) (entities.Page, error) {
	pages, err := repo.Insert(ctx, page)
	if err != nil {
		return entities.Page{}, err
	}

	return pages[0], nil
}
