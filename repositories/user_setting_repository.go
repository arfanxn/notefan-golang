package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/models/entities"

	"github.com/google/uuid"
)

type UserSettingRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewUserSettingRepository(db *sql.DB) *UserSettingRepository {
	return &UserSettingRepository{
		db:          db,
		tableName:   "user_settings",
		columnNames: reflecth.GetFieldJsonTag(entities.UserSetting{}),
	}
}

func (repository *UserSettingRepository) Insert(ctx context.Context, userSettings ...*entities.UserSetting) (
	sql.Result, error) {
	query := buildBatchInsertQuery(repository.tableName, len(userSettings), repository.columnNames...)
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
		errorh.Log(err)
		return result, err
	}
	return result, nil
}

func (repository *UserSettingRepository) Create(ctx context.Context, userSetting *entities.UserSetting) (
	sql.Result, error) {
	return repository.Insert(ctx, userSetting)
}
