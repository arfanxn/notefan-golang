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
	tableName string
	db        *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		tableName: "users",
		db:        db,
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

func (repo *UserRepo) Create(ctx context.Context, user entities.User) (entities.User, error) {
	// Generate uuid for the entity
	user.Id = uuid.New()
	user.CreatedAt = time.Now()

	// Save the entity into database
	query := "INSERT INTO " + repo.tableName + "(id, name, email, password, created_at) VALUES (?, ?, ?, ?)"
	_, err := repo.db.ExecContext(ctx, query, user.Id, user.Name, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		helper.LogIfError(err)
		return entities.User{}, err
	}

	return user, nil
}
