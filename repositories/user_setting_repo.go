package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"

	"github.com/google/uuid"
)

type UserSettingRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewUserSettingRepo(db *sql.DB) *UserSettingRepo {
	return &UserSettingRepo{
		db:          db,
		tableName:   "user_settings",
		columnNames: helper.ReflectGetStructFieldJsonTag(entities.UserSetting{}),
	}
}

func (repo *UserSettingRepo) Insert(ctx context.Context, userSettings ...entities.UserSetting) ([]entities.UserSetting, error) {
	query := buildBatchInsertQuery(repo.tableName, len(userSettings), repo.columnNames...)
	valueArgs := []any{}

	for _, userSetting := range userSettings {
		if userSetting.Id.String() == "" {
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

	stmt, err := repo.db.PrepareContext(ctx, query)
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

func (repo *UserSettingRepo) Create(ctx context.Context, userSetting entities.UserSetting) (entities.UserSetting, error) {
	userSettings, err := repo.Insert(ctx, userSetting)
	if err != nil {
		return entities.UserSetting{}, err
	}

	return userSettings[0], nil
}
