package repositories

import (
	"context"
	"database/sql"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/models/entities"

	"github.com/google/uuid"
)

type RoleRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{
		db:          db,
		tableName:   "roles",
		columnNames: reflecth.GetFieldJsonTag(entities.Role{}),
	}
}

func (repository *RoleRepository) FindByName(ctx context.Context, name string) (entities.Role, error) {
	query := "SELECT id, name FROM " + repository.tableName + " WHERE name = ?"
	var role entities.Role
	rows, err := repository.db.QueryContext(ctx, query, name)
	if err != nil {
		errorh.Log(err)
		return role, err
	}

	if rows.Next() {
		err := rows.Scan(&role.Id, &role.Name)
		if err != nil {
			errorh.Log(err)
			return role, err
		}
	} else {
		return role, exceptions.HTTPNotFound
	}

	return role, nil
}

func (repository *RoleRepository) Insert(ctx context.Context, roles ...entities.Role) ([]entities.Role, error) {
	query := buildBatchInsertQuery(repository.tableName, len(roles), repository.columnNames...)
	valueArgs := []any{}

	for _, role := range roles {
		if role.Id == uuid.Nil {
			role.Id = uuid.New()
		}
		valueArgs = append(valueArgs,
			role.Id,
			role.Name,
		)
	}

	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		errorh.Log(err)
		return roles, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		errorh.Log(err)
		return roles, err
	}
	return roles, nil
}

func (repository *RoleRepository) Create(ctx context.Context, role entities.Role) (entities.Role, error) {
	roles, err := repository.Insert(ctx, role)
	if err != nil {
		return entities.Role{}, err
	}

	return roles[0], nil
}
