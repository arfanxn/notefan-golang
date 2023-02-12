package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"strings"

	"github.com/google/uuid"
)

type PermissionRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewPermissionRepo(db *sql.DB) *PermissionRepo {
	return &PermissionRepo{
		db:          db,
		tableName:   "permissions",
		columnNames: helper.GetStructFieldJsonTag(entities.Permission{}),
	}
}

func (repo *PermissionRepo) Insert(ctx context.Context, permissions ...entities.Permission) ([]entities.Permission, error) {
	query := buildBatchInsertQuery(repo.tableName, len(permissions), repo.columnNames...)
	valueArgs := []any{}

	for _, permission := range permissions {
		if permission.Id.String() == "" {
			permission.Id = uuid.New()
		}
		valueArgs = append(valueArgs, permission.Id, permission.Name)
	}

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
		return permissions, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.LogIfError(err)
		return permissions, err
	}
	return permissions, nil
}

func (repo *PermissionRepo) Create(ctx context.Context, permission entities.Permission) (
	entities.Permission, error) {
	permissions, err := repo.Insert(ctx, permission)
	if err != nil {
		return entities.Permission{}, err
	}

	return permissions[0], nil
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
