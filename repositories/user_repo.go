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

type UserRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db:          db,
		tableName:   "users",
		columnNames: helper.GetStructFieldJsonTag(entities.User{}),
	}
}

func (repo *UserRepo) All(ctx context.Context) ([]entities.User, error) {
	query := "SELECT id, name, email, password, created_at, updated_at FROM " + repo.tableName
	users := []entities.User{}
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
		return users, err
	}

	for rows.Next() {
		user := entities.User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			helper.LogIfError(err)
			return users, err
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return users, exceptions.DataNotFoundError
	}

	return users, nil
}

func (repo *UserRepo) FindByEmail(ctx context.Context, email string) (entities.User, error) {
	query := "SELECT id, name, email, password, created_at, updated_at FROM " + repo.tableName + " WHERE email = ?"
	var user entities.User
	rows, err := repo.db.QueryContext(ctx, query, email)
	if err != nil {
		helper.LogIfError(err)
		return user, err
	}

	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			helper.LogIfError(err)
			return user, err
		}
	} else {
		return user, exceptions.DataNotFoundError
	}

	return user, nil
}

func (repo *UserRepo) Insert(ctx context.Context, users ...entities.User) ([]entities.User, error) {
	query := buildBatchInsertQuery(repo.tableName, len(users), repo.columnNames...)
	valueArgs := []any{}

	for _, user := range users {
		if user.Id.String() == "" {
			user.Id = uuid.New()
		}
		if user.CreatedAt.IsZero() {
			user.CreatedAt = time.Now()
		}
		valueArgs = append(valueArgs,
			user.Id,
			user.Name,
			user.Email,
			user.Password,
			user.CreatedAt,
			user.UpdatedAt,
		)
	}

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
		return users, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.LogIfError(err)
		return users, err
	}
	return users, nil
}

func (repo *UserRepo) Create(ctx context.Context, user entities.User) (entities.User, error) {
	users, err := repo.Insert(ctx, user)
	if err != nil {
		return entities.User{}, err
	}

	return users[0], nil
}
