package repositories

import (
	"context"
	"database/sql"
	"strings"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/models/entities"

	"github.com/google/uuid"
)

type PermissionRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{
		db:          db,
		tableName:   "permissions",
		columnNames: reflecth.GetFieldJsonTag(entities.Permission{}),
	}
}

func (repository *PermissionRepository) Insert(ctx context.Context, permissions ...entities.Permission) ([]entities.Permission, error) {
	query := buildBatchInsertQuery(repository.tableName, len(permissions), repository.columnNames...)
	valueArgs := []any{}

	for _, permission := range permissions {
		if permission.Id == uuid.Nil {
			permission.Id = uuid.New()
		}
		valueArgs = append(valueArgs, permission.Id, permission.Name)
	}

	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		errorh.Log(err)
		return permissions, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		errorh.Log(err)
		return permissions, err
	}
	return permissions, nil
}

func (repository *PermissionRepository) Create(ctx context.Context, permission entities.Permission) (
	entities.Permission, error) {
	permissions, err := repository.Insert(ctx, permission)
	if err != nil {
		return entities.Permission{}, err
	}

	return permissions[0], nil
}

func (repository *PermissionRepository) All(ctx context.Context) (
	[]entities.Permission, error) {
	query := "SELECT id, name FROM " + repository.tableName
	permissions := []entities.Permission{}
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		errorh.Log(err)
		return permissions, err
	}

	for rows.Next() {
		permission := entities.Permission{}
		err := rows.Scan(&permission.Id, &permission.Name)
		if err != nil {
			errorh.Log(err)
			return permissions, err
		}
		permissions = append(permissions, permission)
	}

	if len(permissions) == 0 {
		return permissions, exceptions.HTTPNotFound
	}

	return permissions, nil
}

func (repository *PermissionRepository) FindByNames(ctx context.Context, names ...any) ([]entities.Permission, error) {
	query := "SELECT " + strings.Join(repository.columnNames, ", ") + " FROM " + repository.tableName +
		" WHERE name IN (?" + strings.Repeat(", ?", len(names)-1) + ")"
	permissions := []entities.Permission{}
	rows, err := repository.db.QueryContext(ctx, query, names...)
	if err != nil {
		errorh.Log(err)
		return permissions, err
	}

	for rows.Next() {
		permission := entities.Permission{}
		err := rows.Scan(&permission.Id, &permission.Name)
		if err != nil {
			errorh.Log(err)
			return permissions, err
		}
		permissions = append(permissions, permission)
	}

	if len(permissions) == 0 {
		err := exceptions.HTTPNotFound
		errorh.Log(err)
		return permissions, err
	}

	return permissions, nil
}
