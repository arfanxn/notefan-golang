package repositories

import (
	"context"
	"database/sql"
	"strings"

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

// scanRows scans rows of the database and returns it as structs, and returns error if any error has occurred.
func (repository *PermissionRepository) scanRows(rows *sql.Rows) (permissions []entities.Permission, err error) {
	for rows.Next() {
		var permission entities.Permission
		err := rows.Scan(&permission.Id, &permission.Name)
		if err != nil {
			return permissions, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

// scanRow scans only a row of the database and returns it as struct, and returns error if any error has occurred.
func (repository *PermissionRepository) scanRow(rows *sql.Rows) (entities.Permission, error) {
	permissions, err := repository.scanRows(rows)
	if len(permissions) == 0 {
		return entities.Permission{}, err
	}
	return permissions[0], nil
}

func (repository *PermissionRepository) All(ctx context.Context) (
	permissions []entities.Permission, err error) {
	query := "SELECT id, name FROM " + repository.tableName
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		return permissions, err
	}
	permissions, err = repository.scanRows(rows)
	return permissions, err
}

// GetByNames retrieves data from database table by names
func (repository *PermissionRepository) GetByNames(ctx context.Context, names ...any) (
	permissions []entities.Permission, err error) {
	query := "SELECT " + strings.Join(repository.columnNames, ", ") + " FROM " + repository.tableName +
		" WHERE name IN (?" + strings.Repeat(", ?", len(names)-1) + ")"
	rows, err := repository.db.QueryContext(ctx, query, names...)
	if err != nil {
		return
	}
	permissions, err = repository.scanRows(rows)
	return
}

func (repository *PermissionRepository) Insert(ctx context.Context, permissions ...*entities.Permission) (
	sql.Result, error) {
	query := buildBatchInsertQuery(repository.tableName, len(permissions), repository.columnNames...)
	valueArgs := []any{}
	for _, permission := range permissions {
		if permission.Id == uuid.Nil {
			permission.Id = uuid.New()
		}
		valueArgs = append(valueArgs, permission.Id, permission.Name)
	}
	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (repository *PermissionRepository) Create(ctx context.Context, permission *entities.Permission) (
	sql.Result, error) {
	return repository.Insert(ctx, permission)
}
