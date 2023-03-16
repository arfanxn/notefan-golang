package repositories

import (
	"context"
	"database/sql"

	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/query_reqs"

	"github.com/google/uuid"
)

type RoleRepository struct {
	db     *sql.DB
	Query  query_reqs.Query
	entity entities.Role
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{
		db:     db,
		Query:  query_reqs.Default(),
		entity: entities.Role{},
	}
}

// All retrives all roles
func (repository *RoleRepository) All(ctx context.Context) (
	roles []entities.Role, err error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.entity.GetColumnNames()) +
		" FROM " + repository.entity.GetTableName()
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		return
	}
	for rows.Next() {
		var role entities.Role
		rows.Scan(&role.Id, &role.Name)
		roles = append(roles, role)
	}
	return
}

// FindByName finds by name
func (repository *RoleRepository) FindByName(ctx context.Context, name string) (role entities.Role, err error) {
	query := "SELECT id, name FROM " + repository.entity.GetTableName() +
		" WHERE name = ?"
	rows, err := repository.db.QueryContext(ctx, query, name)
	if err != nil {
		return
	}
	if rows.Next() {
		err = rows.Scan(&role.Id, &role.Name)
		if err != nil {
			return
		}
	}
	return role, err
}

// Insert inserts data into database table
func (repository *RoleRepository) Insert(ctx context.Context, roles ...*entities.Role) (sql.Result, error) {
	query := buildBatchInsertQuery(
		repository.entity.GetTableName(),
		len(roles),
		repository.entity.GetColumnNames()...,
	)
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
	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return result, err
	}
	return result, nil
}

// Create creates data into database table
func (repository *RoleRepository) Create(ctx context.Context, role *entities.Role) (sql.Result, error) {
	return repository.Insert(ctx, role)
}
