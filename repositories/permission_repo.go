package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"strings"
)

type PermissionRepo struct {
	tableName string
	db        *sql.DB
}

func NewPermissionRepo(db *sql.DB) *PermissionRepo {
	return &PermissionRepo{
		tableName: "permissions",
		db:        db,
	}
}

func (repo *PermissionRepo) All(ctx context.Context) (
	[]entities.Permission, error) {
	query := "SELECT id, name FROM " + repo.tableName
	permissions := []entities.Permission{}
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
		return permissions, err
	}

	for rows.Next() {
		permission := entities.Permission{}
		err := rows.Scan(&permission.Id, &permission.Name)
		if err != nil {
			helper.LogIfError(err)
			return permissions, err
		}
		permissions = append(permissions, permission)
	}

	if len(permissions) == 0 {
		return permissions, exceptions.DataNotFoundError
	}

	return permissions, nil
}

func (repo *PermissionRepo) GetByNames(ctx context.Context, names ...any) (
	[]entities.Permission, error) {
	query := "SELECT id, name FROM " + repo.tableName +
		" WHERE name IN (?" + strings.Repeat(", ?", len(names)-1) + ")"
	permissions := []entities.Permission{}
	rows, err := repo.db.QueryContext(ctx, query, names...)
	if err != nil {
		helper.LogIfError(err)
		return permissions, err
	}

	for rows.Next() {
		permission := entities.Permission{}
		err := rows.Scan(&permission.Id, &permission.Name)
		if err != nil {
			helper.LogIfError(err)
			return permissions, err
		}
		permissions = append(permissions, permission)
	}

	if len(permissions) == 0 {
		err := exceptions.DataNotFoundError
		helper.LogIfError(err)
		return permissions, err
	}

	return permissions, nil
}

func (repo *PermissionRepo) FindByName(ctx context.Context, name string) (entities.Permission, error) {
	permissions, err := repo.GetByNames(ctx, name)
	if err != nil {
		return entities.Permission{}, err
	}
	return permissions[0], nil
}
