package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/entities"

	"github.com/google/uuid"
)

type NotificationRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{
		db:          db,
		tableName:   "notifications",
		columnNames: helper.ReflectGetStructFieldJsonTag(entities.Notification{}),
	}
}

func (repository *NotificationRepository) All(ctx context.Context) ([]entities.Notification, error) {
	query := "SELECT " + helper.DBSliceColumnsToStr(repository.columnNames) + " FROM " + repository.tableName
	notifications := []entities.Notification{}
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		helper.ErrorLog(err)
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
			helper.ErrorLog(err)
			return notifications, err
		}
		notifications = append(notifications, notification)
	}

	if len(notifications) == 0 {
		return notifications, exceptions.HTTPNotFound
	}

	return notifications, nil
}

func (repository *NotificationRepository) Insert(ctx context.Context, notifications ...entities.Notification) ([]entities.Notification, error) {
	query := buildBatchInsertQuery(repository.tableName, len(notifications), repository.columnNames...)
	valueArgs := []any{}

	for _, notification := range notifications {
		if notification.Id == uuid.Nil {
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

	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		helper.ErrorLog(err)
		return notifications, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.ErrorLog(err)
		return notifications, err
	}
	return notifications, nil
}

func (repository *NotificationRepository) Create(ctx context.Context, notification entities.Notification) (entities.Notification, error) {
	notifications, err := repository.Insert(ctx, notification)
	if err != nil {
		return entities.Notification{}, err
	}

	return notifications[0], nil
}
