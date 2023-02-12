package repositories

import "database/sql"

type NotificationUserRepo struct {
	tableName string
	db        *sql.DB
}

func NewNotificationUserRepo(db *sql.DB) *NotificationUserRepo {
	return &NotificationUserRepo{
		tableName: "notification_user",
		db:        db,
	}
}
