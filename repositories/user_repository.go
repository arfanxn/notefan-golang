package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"

	"github.com/google/uuid"
)

// UserRepository represents a repository for user model/entity
type UserRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

// Instantiate a UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db:          db,
		tableName:   "users",
		columnNames: reflecth.GetFieldJsonTag(entities.User{}),
	}
}

// scanRows scans rows of the database and returns it as structs, and returns error if any error has occurred.
func (repository *UserRepository) scanRows(rows *sql.Rows) ([]entities.User, error) {
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
		errorh.Panic(err) // panic if scan fails
		users = append(users, user)
	}

	if len(users) == 0 {
		return users, exceptions.HTTPNotFound
	}
	return users, nil
}

// scanRow scans only a row of the database and returns it as struct, and returns error if any error has occurred.
func (repository *UserRepository) scanRow(rows *sql.Rows) (entities.User, error) {
	users, err := repository.scanRows(rows)
	if err != nil {
		return entities.User{}, err
	}
	return users[0], nil
}

// All retrives all users
func (repository *UserRepository) All(ctx context.Context) ([]entities.User, error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.columnNames) + " FROM " + repository.tableName
	rows, err := repository.db.QueryContext(ctx, query)
	errorh.Panic(err) // panic if query error
	return repository.scanRows(rows)
}

// Find finds a user by id
func (repository *UserRepository) Find(ctx context.Context, id string) (entities.User, error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.columnNames) + " FROM " + repository.tableName +
		" WHERE id = ?"
	rows, err := repository.db.QueryContext(ctx, query, id)
	errorh.Panic(err) // panic if query error

	return repository.scanRow(rows)
}

// FindByEmail finds a user by email address
func (repository *UserRepository) FindByEmail(ctx context.Context, email string) (entities.User, error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.columnNames) + " FROM " + repository.tableName +
		" WHERE email = ?"
	rows, err := repository.db.QueryContext(ctx, query, email)
	errorh.Panic(err) // panic if query error

	return repository.scanRow(rows)
}

// Insert inserts users into the database
func (repository *UserRepository) Insert(ctx context.Context, users ...entities.User) ([]entities.User, error) {
	query := buildBatchInsertQuery(repository.tableName, len(users), repository.columnNames...)
	valueArgs := []any{}

	for _, user := range users {
		if user.Id == uuid.Nil {
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

	stmt, err := repository.db.PrepareContext(ctx, query)
	errorh.Panic(err) // panic if query error
	_, err = stmt.ExecContext(ctx, valueArgs...)
	errorh.Panic(err) // panic if query error

	return users, nil
}

// Insert inserts a user into the database
func (repository *UserRepository) Create(ctx context.Context, user entities.User) (entities.User, error) {
	users, err := repository.Insert(ctx, user)
	if err != nil {
		return entities.User{}, err
	}

	return users[0], nil
}
