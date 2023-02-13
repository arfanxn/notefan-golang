package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
)

type PermissionRoleRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewPermissionRoleRepo(db *sql.DB) *PermissionRoleRepo {
	return &PermissionRoleRepo{
		db:          db,
		tableName:   "permission_role",
		columnNames: helper.ReflectGetStructFieldJsonTag(entities.PermissionRole{}),
	}
}

func (repo *PermissionRoleRepo) Insert(ctx context.Context, permissionRoles ...entities.PermissionRole) ([]entities.PermissionRole, error) {
	query := buildBatchInsertQuery(repo.tableName, len(permissionRoles), repo.columnNames...)
	valueArgs := []any{}

	for _, permissionRole := range permissionRoles {
		valueArgs = append(valueArgs,
			permissionRole.PermissionId, permissionRole.RoleId, permissionRole.CreatedAt)
	}

	stmt, err := repo.db.PrepareContext(ctx, query)
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

func (repo *PermissionRoleRepo) Create(ctx context.Context, permissionRole entities.PermissionRole) (
	entities.PermissionRole, error) {
	permissionRoles, err := repo.Insert(ctx, permissionRole)
	if err != nil {
		return entities.PermissionRole{}, err
	}

	return permissionRoles[0], nil
}
