package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/helper"
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
		columnNames: helper.ReflectGetStructFieldJsonTag(entities.UserSetting{}),
	}
}

func (repository *UserSettingRepository) Insert(ctx context.Context, userSettings ...entities.UserSetting) ([]entities.UserSetting, error) {
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

	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		helper.ErrorLog(err)
		return userSettings, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.ErrorLog(err)
		return userSettings, err
	}
	return userSettings, nil
}

func (repository *UserSettingRepository) Create(ctx context.Context, userSetting entities.UserSetting) (entities.UserSetting, error) {
	userSettings, err := repository.Insert(ctx, userSetting)
	if err != nil {
		return entities.UserSetting{}, err
	}

	return userSettings[0], nil
}
