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

// UserRepo represents a repository for user model/entity
type UserRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

// Instantiate a UserRepo
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db:          db,
		tableName:   "users",
		columnNames: helper.ReflectGetStructFieldJsonTag(entities.User{}),
	}
}

// scanRows scans rows of the database and returns it as structs, and returns error if any error has occurred.
func (repo *UserRepo) scanRows(rows *sql.Rows) ([]entities.User, error) {
	users := []entities.User{}
	for rows.Next() {
		user := entities.User{}
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		helper.ErrorPanic(err) // panic if scan fails
		users = append(users, user)
	}

	if len(users) == 0 {
		return users, exceptions.HTTPNotFound
	}
	return users, nil
}

// scanRow scans only a row of the database and returns it as struct, and returns error if any error has occurred.
func (repo *UserRepo) scanRow(rows *sql.Rows) (entities.User, error) {
	users, err := repo.scanRows(rows)
	if err != nil {
		return entities.User{}, err
	}
	return users[0], nil
}

// All retrives all users
func (repo *UserRepo) All(ctx context.Context) ([]entities.User, error) {
	query := "SELECT " + helper.DBSliceColumnsToStr(repo.columnNames) + " FROM " + repo.tableName
	rows, err := repo.db.QueryContext(ctx, query)
	helper.ErrorPanic(err) // panic if query error
	return repo.scanRows(rows)
}

// FindByEmail finds a user by email address
func (repo *UserRepo) FindByEmail(ctx context.Context, email string) (entities.User, error) {
	query := "SELECT " + helper.DBSliceColumnsToStr(repo.columnNames) + " FROM " + repo.tableName +
		" WHERE email = ?"
	rows, err := repo.db.QueryContext(ctx, query, email)
	helper.ErrorPanic(err) // panic if query error

	return repo.scanRow(rows)
}

// Insert inserts users into the database
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
	helper.ErrorPanic(err) // panic if query error
	_, err = stmt.ExecContext(ctx, valueArgs...)
	helper.ErrorPanic(err) // panic if query error

	return users, nil
}

// Insert inserts a user into the database
func (repo *UserRepo) Create(ctx context.Context, user entities.User) (entities.User, error) {
	users, err := repo.Insert(ctx, user)
	if err != nil {
		return entities.User{}, err
	}

	return users[0], nil
}
