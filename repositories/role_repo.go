package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
)

type RoleRepo struct {
	tableName string
	db        *sql.DB
}

func NewRoleRepo(db *sql.DB) *RoleRepo {
	return &RoleRepo{
		tableName: "roles",
		db:        db,
	}
}

func (repo *RoleRepo) Get() {

}

func (repo *RoleRepo) FindByName(ctx context.Context, name string) (entities.Role, error) {
	query := "SELECT id, name FROM " + repo.tableName + " WHERE name = ?"
	var role entities.Role
	rows, err := repo.db.QueryContext(ctx, query, name)
	if err != nil {
		helper.LogIfError(err)
		return role, err
	}

	if rows.Next() {
		err := rows.Scan(&role.Id, &role.Name)
		if err != nil {
			helper.LogIfError(err)
			return role, err
		}
	} else {
		return role, exceptions.DataNotFoundError
	}

	return role, nil
}
