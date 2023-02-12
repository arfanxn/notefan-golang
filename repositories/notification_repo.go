package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"

	"github.com/google/uuid"
)

type NotificationRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewNotificationRepo(db *sql.DB) *NotificationRepo {
	return &NotificationRepo{
		db:          db,
		tableName:   "notifications",
		columnNames: helper.GetStructFieldJsonTag(entities.Notification{}),
	}
}

func (repo *NotificationRepo) All(ctx context.Context) ([]entities.Notification, error) {
	query := "SELECT id, name, description, domain, created_at, updated_at FROM " + repo.tableName
	notifications := []entities.Notification{}
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
		return notifications, err
	}

	for rows.Next() {
		notification := entities.Notification{}
		err := rows.Scan(
			&notification.Id,
			&notification.ObjectType,
			&notification.ObjectId,
			&notification.Title,
			&notification.Type,
			&notification.Body,
			&notification.ArchivedAt,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			helper.LogIfError(err)
			return notifications, err
		}
		notifications = append(notifications, notification)
	}

	if len(notifications) == 0 {
		return notifications, exceptions.DataNotFoundError
	}

	return notifications, nil
}

func (repo *NotificationRepo) Insert(ctx context.Context, notifications ...entities.Notification) ([]entities.Notification, error) {
	query := buildBatchInsertQuery(repo.tableName, len(notifications), repo.columnNames...)
	valueArgs := []any{}

	for _, notification := range notifications {
		if notification.Id.String() == "" {
			notification.Id = uuid.New()
		}
		if notification.CreatedAt.IsZero() {
			notification.CreatedAt = time.Now()
		}
		valueArgs = append(valueArgs,
			notification.Id,
			notification.ObjectType,
			notification.ObjectId,
			notification.Title,
			notification.Type,
			notification.Body,
			notification.ArchivedAt,
			notification.CreatedAt,
			notification.UpdatedAt,
		)
	}

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
		return notifications, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.LogIfError(err)
		return notifications, err
	}
	return notifications, nil
}

func (repo *NotificationRepo) Create(ctx context.Context, notification entities.Notification) (entities.Notification, error) {
	notifications, err := repo.Insert(ctx, notification)
	if err != nil {
		return entities.Notification{}, err
	}

	return notifications[0], nil
}
