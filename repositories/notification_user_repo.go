package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
)

type NotificationUserRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewNotificationUserRepo(db *sql.DB) *NotificationUserRepo {
	return &NotificationUserRepo{
		db:          db,
		tableName:   "notification_user",
		columnNames: helper.GetStructFieldJsonTag(entities.NotificationUser{}),
	}
}

func (repo *NotificationUserRepo) All(ctx context.Context) ([]entities.NotificationUser, error) {
	query := "SELECT " + helper.DBSliceColumnsToStr(repo.columnNames) + " FROM " + repo.tableName
	notificationUsers := []entities.NotificationUser{}
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
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
			helper.LogIfError(err)
			return notificationUsers, err
		}
		notificationUsers = append(notificationUsers, notificationUser)
	}

	if len(notificationUsers) == 0 {
		return notificationUsers, exceptions.DataNotFoundError
	}

	return notificationUsers, nil
}

func (repo *NotificationUserRepo) Insert(ctx context.Context, notificationUsers ...entities.NotificationUser) ([]entities.NotificationUser, error) {
	query := buildBatchInsertQuery(repo.tableName, len(notificationUsers), repo.columnNames...)
	valueArgs := []any{}

	for _, notificationUser := range notificationUsers {
		valueArgs = append(valueArgs,
			notificationUser.NotificationId,
			notificationUser.NotifierId,
			notificationUser.NotifiedId,
		)
	}

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
		return notificationUsers, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.LogIfError(err)
		return notificationUsers, err
	}
	return notificationUsers, nil
}

func (repo *NotificationUserRepo) Create(ctx context.Context, notificationUser entities.NotificationUser) (entities.NotificationUser, error) {
	notificationUsers, err := repo.Insert(ctx, notificationUser)
	if err != nil {
		return entities.NotificationUser{}, err
	}

	return notificationUsers[0], nil
}
