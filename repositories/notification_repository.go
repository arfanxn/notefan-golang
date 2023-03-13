package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/query_reqs"

	"github.com/google/uuid"
)

type NotificationRepository struct {
	db     *sql.DB
	Query  query_reqs.Query
	entity entities.Notification
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{
		db:     db,
		Query:  query_reqs.Default(),
		entity: entities.Notification{},
	}
}

func (repository *NotificationRepository) All(ctx context.Context) ([]entities.Notification, error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.entity.GetColumnNames()) +
		" FROM " + repository.entity.GetTableName()
	notifications := []entities.Notification{}
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
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
			return notifications, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func (repository *NotificationRepository) Insert(ctx context.Context, notifications ...*entities.Notification) (sql.Result, error) {
	query := buildBatchInsertQuery(
		repository.entity.GetTableName(),
		len(notifications),
		repository.entity.GetColumnNames()...,
	)
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

	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (repository *NotificationRepository) Create(ctx context.Context, notification *entities.Notification) (
	sql.Result, error) {
	return repository.Insert(ctx, notification)
}
