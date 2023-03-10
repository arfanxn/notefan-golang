package repositories

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/query_reqs"

	"github.com/google/uuid"
)

// UserRepository represents a repository for user model/entity
type UserRepository struct {
	db        *sql.DB
	Query     query_reqs.Query
	entity    entities.User
	mutex     sync.Mutex
	waitGroup *sync.WaitGroup
}

// Instantiate a UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db:        db,
		Query:     query_reqs.Default(),
		entity:    entities.User{},
		mutex:     sync.Mutex{},
		waitGroup: new(sync.WaitGroup),
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
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// scanRow scans only a row of the database and returns it as struct, and returns error if any error has occurred.
func (repository *UserRepository) scanRow(rows *sql.Rows) (user entities.User, err error) {
	users, err := repository.scanRows(rows)
	if err != nil {
		return
	}
	if len(users) == 0 {
		return
	}
	user, err = users[0], nil
	return
}

// All retrives all users
func (repository *UserRepository) All(ctx context.Context) (
	users []entities.User, err error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.entity.GetColumnNames()) +
		" FROM " + repository.entity.GetTableName()
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		return users, err
	}
	users, err = repository.scanRows(rows)
	return
}

// Find finds a user by id
func (repository *UserRepository) Find(ctx context.Context, id string) (user entities.User, err error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.entity.GetColumnNames()) +
		" FROM " + repository.entity.GetTableName() +
		" WHERE id = ?"
	rows, err := repository.db.QueryContext(ctx, query, id)
	if err != nil {
		return user, err
	}

	user, err = repository.scanRow(rows)
	if err != nil {
		return user, err
	}
	return user, nil
}

// FindByEmail finds a user by email address
func (repository *UserRepository) FindByEmail(ctx context.Context, email string) (user entities.User, err error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.entity.GetColumnNames()) +
		" FROM " + repository.entity.GetTableName() +
		" WHERE email = ?"
	rows, err := repository.db.QueryContext(ctx, query, email)
	if err != nil {
		return user, err
	}

	user, err = repository.scanRow(rows)
	if err != nil {
		return user, err
	}
	return user, nil
}

// Insert inserts users into the database
func (repository *UserRepository) Insert(ctx context.Context, users ...*entities.User) (sql.Result, error) {
	query := buildBatchInsertQuery(
		repository.entity.GetTableName(),
		len(users),
		repository.entity.GetColumnNames()...,
	)
	valueArgs := []any{}

	for _, user := range users {
		repository.waitGroup.Add(1)

		go func(wg *sync.WaitGroup, user *entities.User) {

			defer wg.Done()

			repository.mutex.Lock()
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
			repository.mutex.Unlock()
		}(repository.waitGroup, user)
	}

	repository.waitGroup.Wait()

	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return result, err
	}
	return result, nil
}

// Create creates user into database
func (repository *UserRepository) Create(ctx context.Context, user *entities.User) (sql.Result, error) {
	result, err := repository.Insert(ctx, user)
	if err != nil {
		return result, err
	}

	return result, nil
}

// UpdateById
func (repository *UserRepository) UpdateById(ctx context.Context, user *entities.User) (sql.Result, error) {
	query := buildUpdateQuery(repository.entity.GetTableName(), repository.entity.GetColumnNames()...) +
		" WHERE id = ?"

	// Refresh entity updated at
	user.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	result, err := repository.db.ExecContext(ctx, query,
		user.Id, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt, user.Id)

	return result, err
}
