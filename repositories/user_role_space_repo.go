package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"
)

type UserRoleSpaceRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewUserRoleSpaceRepo(db *sql.DB) *UserRoleSpaceRepo {
	return &UserRoleSpaceRepo{
		db:          db,
		tableName:   "user_role_space",
		columnNames: helper.ReflectGetStructFieldJsonTag(entities.UserRoleSpace{}),
	}
}

func (repo *UserRoleSpaceRepo) Insert(ctx context.Context, userRoleSpaces ...entities.UserRoleSpace) (
	[]entities.UserRoleSpace, error) {
	query := buildBatchInsertQuery(repo.tableName, len(userRoleSpaces), repo.columnNames...)
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

	stmt, err := repo.db.PrepareContext(ctx, query)
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

func (repo *UserRoleSpaceRepo) Create(ctx context.Context, userRoleSpace entities.UserRoleSpace) (
	entities.UserRoleSpace, error) {
	userRoleSpaces, err := repo.Insert(ctx, userRoleSpace)
	if err != nil {
		return entities.UserRoleSpace{}, err
	}

	return userRoleSpaces[0], nil
}
