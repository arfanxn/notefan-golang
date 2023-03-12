package repositories

import (
	"context"
	"database/sql"

	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/models/entities"
)

type PermissionRoleRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewPermissionRoleRepository(db *sql.DB) *PermissionRoleRepository {
	return &PermissionRoleRepository{
		db:          db,
		tableName:   "permission_role",
		columnNames: reflecth.GetFieldJsonTag(entities.PermissionRole{}),
	}
}

// Insert inserts into the database table
func (repository *PermissionRoleRepository) Insert(
	ctx context.Context, permissionRoles ...*entities.PermissionRole) (
	sql.Result, error) {
	query := buildBatchInsertQuery(repository.tableName, len(permissionRoles), repository.columnNames...)
	var valueArgs []any
	for _, permissionRole := range permissionRoles {
		valueArgs = append(valueArgs,
			permissionRole.PermissionId, permissionRole.RoleId, permissionRole.CreatedAt)
	}
	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return result, err
	}
	return result, nil
}

// Create do same as Insert but singularly
func (repository *PermissionRoleRepository) Create(ctx context.Context, permissionRole *entities.PermissionRole) (
	sql.Result, error) {
	return repository.Insert(ctx, permissionRole)
}
