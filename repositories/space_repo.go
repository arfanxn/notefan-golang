package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
)

type SpaceRepo struct {
	tableName string
	db        *sql.DB
}

func NewSpaceRepo(db *sql.DB) *SpaceRepo {
	return &SpaceRepo{
		tableName: "spaces",
		db:        db,
	}
}

func (repo *SpaceRepo) All(ctx context.Context) ([]entities.Space, error) {
	query := "SELECT id, name, description, domain, created_at, updated_at FROM " + repo.tableName
	spaces := []entities.Space{}
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
		return spaces, err
	}

	for rows.Next() {
		space := entities.Space{}
		err := rows.Scan(
			&space.Id, &space.Name, &space.Description,
			&space.Domain, &space.CreatedAt, &space.UpdatedAt,
		)
		if err != nil {
			helper.LogIfError(err)
			return spaces, err
		}
		spaces = append(spaces, space)
	}

	if len(spaces) == 0 {
		return spaces, exceptions.DataNotFoundError
	}

	return spaces, nil
}
