package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/helper"
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
		columnNames: helper.ReflectGetStructFieldJsonTag(entities.UserRoleSpace{}),
	}
}

func (repository *UserRoleSpaceRepository) Insert(ctx context.Context, userRoleSpaces ...entities.UserRoleSpace) (
	[]entities.UserRoleSpace, error) {
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

	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		helper.ErrorLog(err)
		return userRoleSpaces, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.ErrorLog(err)
		return userRoleSpaces, err
	}
	return userRoleSpaces, nil
}

func (repository *UserRoleSpaceRepository) Create(ctx context.Context, userRoleSpace entities.UserRoleSpace) (
	entities.UserRoleSpace, error) {
	userRoleSpaces, err := repository.Insert(ctx, userRoleSpace)
	if err != nil {
		return entities.UserRoleSpace{}, err
	}

	return userRoleSpaces[0], nil
}
