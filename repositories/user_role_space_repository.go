package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/models/entities"
)

type UserRoleSpaceRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewUserRoleSpaceRepository(db *sql.DB) *UserRoleSpaceRepository {
	return &UserRoleSpaceRepository{
		db:          db,
		tableName:   "user_role_space",
		columnNames: reflecth.GetFieldJsonTag(entities.UserRoleSpace{}),
	}
}

func (repository *UserRoleSpaceRepository) Insert(ctx context.Context, userRoleSpaces ...*entities.UserRoleSpace) (
	sql.Result, error) {
	query := buildBatchInsertQuery(repository.tableName, len(userRoleSpaces), repository.columnNames...)
	valueArgs := []any{}

	for _, userRoleSpace := range userRoleSpaces {
		if userRoleSpace.CreatedAt.IsZero() {
			userRoleSpace.CreatedAt = time.Now()
		}
		valueArgs = append(valueArgs,
			userRoleSpace.UserId,
			userRoleSpace.RoleId,
			userRoleSpace.SpaceId,
			userRoleSpace.CreatedAt,
			userRoleSpace.UpdatedAt,
		)
	}

	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	return result, err
}

func (repository *UserRoleSpaceRepository) Create(ctx context.Context, userRoleSpace *entities.UserRoleSpace) (
	sql.Result, error) {
	result, err := repository.Insert(ctx, userRoleSpace)
	if err != nil {
		return result, err
	}

	return result, nil
}
