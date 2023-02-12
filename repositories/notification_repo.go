package repositories

import "database/sql"

type NotificationRepo struct {
	tableName string
	db        *sql.DB
}

func NewNotificationRepo(db *sql.DB) *NotificationRepo {
	return &NotificationRepo{
		tableName: "notifications",
		db:        db,
	}
}
