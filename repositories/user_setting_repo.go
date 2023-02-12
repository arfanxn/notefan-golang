package repositories

import "database/sql"

type UserSettingRepo struct {
	tableName string
	db        *sql.DB
}

func NewUserSettingRepo(db *sql.DB) *UserSettingRepo {
	return &UserSettingRepo{
		tableName: "user_settings",
		db:        db,
	}
}
