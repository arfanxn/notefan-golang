package repositories

import "database/sql"

type PermissionRoleRepo struct {
	tableName string
	db        *sql.DB
}

func NewPermissionRoleRepo(db *sql.DB) *PermissionRoleRepo {
	return &PermissionRoleRepo{
		tableName: "permission_role",
		db:        db,
	}
}
