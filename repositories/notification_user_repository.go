package repositories

import (
	"context"
	"database/sql"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"
)

type NotificationUserRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewNotificationUserRepository(db *sql.DB) *NotificationUserRepository {
	return &NotificationUserRepository{
		db:          db,
		tableName:   "notification_user",
		columnNames: reflecth.GetFieldJsonTag(entities.NotificationUser{}),
	}
}

func (repository *NotificationUserRepository) All(ctx context.Context) ([]entities.NotificationUser, error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.columnNames) + " FROM " + repository.tableName
	notificationUsers := []entities.NotificationUser{}
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		errorh.Log(err)
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
			errorh.Log(err)
			return notificationUsers, err
		}
		notificationUsers = append(notificationUsers, notificationUser)
	}

	if len(notificationUsers) == 0 {
		return notificationUsers, exceptions.HTTPNotFound
	}

	return notificationUsers, nil
}

func (repository *NotificationUserRepository) Insert(ctx context.Context, notificationUsers ...entities.NotificationUser) ([]entities.NotificationUser, error) {
	query := buildBatchInsertQuery(repository.tableName, len(notificationUsers), repository.columnNames...)
	valueArgs := []any{}

	for _, notificationUser := range notificationUsers {
		valueArgs = append(valueArgs,
			notificationUser.NotificationId,
			notificationUser.NotifierId,
			notificationUser.NotifiedId,
		)
	}

	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		errorh.Log(err)
		return notificationUsers, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		errorh.Log(err)
		return notificationUsers, err
	}
	return notificationUsers, nil
}

func (repository *NotificationUserRepository) Create(ctx context.Context, notificationUser entities.NotificationUser) (entities.NotificationUser, error) {
	notificationUsers, err := repository.Insert(ctx, notificationUser)
	if err != nil {
		return entities.NotificationUser{}, err
	}

	return notificationUsers[0], nil
}
