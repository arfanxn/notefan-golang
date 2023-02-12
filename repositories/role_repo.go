package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/entities"

	"github.com/google/uuid"
)

type RoleRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewRoleRepo(db *sql.DB) *RoleRepo {
	return &RoleRepo{
		db:          db,
		tableName:   "roles",
		columnNames: helper.GetStructFieldJsonTag(entities.Role{}),
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

func (repo *RoleRepo) Insert(ctx context.Context, roles ...entities.Role) ([]entities.Role, error) {
	query := buildBatchInsertQuery(repo.tableName, len(roles), repo.columnNames...)
	valueArgs := []any{}

	for _, role := range roles {
		if role.Id.String() == "" {
			role.Id = uuid.New()
		}
		valueArgs = append(valueArgs,
			role.Id,
			role.Name,
		)
	}

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
		return roles, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.LogIfError(err)
		return roles, err
	}
	return roles, nil
}

func (repo *RoleRepo) Create(ctx context.Context, role entities.Role) (entities.Role, error) {
	roles, err := repo.Insert(ctx, role)
	if err != nil {
		return entities.Role{}, err
	}

	return roles[0], nil
}
