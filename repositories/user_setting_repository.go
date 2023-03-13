package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/query_reqs"

	"github.com/google/uuid"
)

type UserSettingRepository struct {
	db     *sql.DB
	Query  query_reqs.Query
	entity entities.UserSetting
}

func NewUserSettingRepository(db *sql.DB) *UserSettingRepository {
	return &UserSettingRepository{
		db:     db,
		Query:  query_reqs.Default(),
		entity: entities.UserSetting{},
	}
}

func (repository *UserSettingRepository) Insert(ctx context.Context, userSettings ...*entities.UserSetting) (
	sql.Result, error) {
	query := buildBatchInsertQuery(
		repository.entity.GetTableName(),
		len(userSettings),
		repository.entity.GetColumnNames()...,
	)
	valueArgs := []any{}

	for _, userSetting := range userSettings {
		if userSetting.Id == uuid.Nil {
			userSetting.Id = uuid.New()
		}
		if userSetting.CreatedAt.IsZero() {
			userSetting.CreatedAt = time.Now()
		}
		valueArgs = append(valueArgs,
			userSetting.Id,
			userSetting.UserId,
			userSetting.Key,
			userSetting.Value,
			userSetting.CreatedAt,
			userSetting.UpdatedAt,
		)
	}

	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (repository *UserSettingRepository) Create(ctx context.Context, userSetting *entities.UserSetting) (
	sql.Result, error) {
	return repository.Insert(ctx, userSetting)
}
