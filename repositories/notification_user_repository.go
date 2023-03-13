package repositories

import (
	"context"
	"database/sql"

	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/query_reqs"
)

type NotificationUserRepository struct {
	db     *sql.DB
	Query  query_reqs.Query
	entity entities.NotificationUser
}

func NewNotificationUserRepository(db *sql.DB) *NotificationUserRepository {
	return &NotificationUserRepository{
		db:     db,
		Query:  query_reqs.Default(),
		entity: entities.NotificationUser{},
	}
}

func (repository *NotificationUserRepository) All(ctx context.Context) ([]entities.NotificationUser, error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.entity.GetColumnNames()) +
		" FROM " + repository.entity.GetTableName()
	notificationUsers := []entities.NotificationUser{}
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		return notificationUsers, err
	}
	for rows.Next() {
		notificationUser := entities.NotificationUser{}
		err := rows.Scan(
			&notificationUser.NotificationId,
			&notificationUser.NotifierId,
			&notificationUser.NotifiedId,
		)
		if err != nil {
			return notificationUsers, err
		}
		notificationUsers = append(notificationUsers, notificationUser)
	}
	return notificationUsers, nil
}

func (repository *NotificationUserRepository) Insert(ctx context.Context, notificationUsers ...*entities.NotificationUser) (sql.Result, error) {
	query := buildBatchInsertQuery(
		repository.entity.GetTableName(),
		len(notificationUsers),
		repository.entity.GetColumnNames()...,
	)
	valueArgs := []any{}
	for _, notificationUser := range notificationUsers {
		valueArgs = append(valueArgs,
			notificationUser.NotificationId,
			notificationUser.NotifierId,
			notificationUser.NotifiedId,
		)
	}
	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (repository *NotificationUserRepository) Create(ctx context.Context, notificationUser *entities.NotificationUser) (sql.Result, error) {
	return repository.Insert(ctx, notificationUser)
}
