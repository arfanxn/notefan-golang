package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/helper"
	"notefan-golang/models/entities"

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

func (repo *UserRepo) Create(ctx context.Context, user entities.User) (entities.User, error) {
	// Generate uuid for the entity
	user.Id = uuid.New()

	// Save the entity into database
	query := "INSERT INTO " + repo.tableName + "(id, name, email, password) VALUES (?, ?, ?, ?)"
	_, err := repo.db.ExecContext(ctx, query, user.Id, user.Name, user.Email, user.Password)
	if err != nil {
		helper.LogIfError(err)
		return entities.User{}, err
	}

	return user, err
}
