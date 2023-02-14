package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
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
		columnNames: helper.ReflectGetStructFieldJsonTag(entities.PermissionRole{}),
	}
}

func (repository *PermissionRoleRepository) Insert(ctx context.Context, permissionRoles ...entities.PermissionRole) ([]entities.PermissionRole, error) {
	query := buildBatchInsertQuery(repository.tableName, len(permissionRoles), repository.columnNames...)
	valueArgs := []any{}

	for _, permissionRole := range permissionRoles {
		valueArgs = append(valueArgs,
			permissionRole.PermissionId, permissionRole.RoleId, permissionRole.CreatedAt)
	}

	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		helper.ErrorLog(err)
		return permissionRoles, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.ErrorLog(err)
		return permissionRoles, err
	}
	return permissionRoles, nil
}

func (repository *PermissionRoleRepository) Create(ctx context.Context, permissionRole entities.PermissionRole) (
	entities.PermissionRole, error) {
	permissionRoles, err := repository.Insert(ctx, permissionRole)
	if err != nil {
		return entities.PermissionRole{}, err
	}

	return permissionRoles[0], nil
}
