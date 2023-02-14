package seeders

import (
	"context"
	"database/sql"
	"notefan-golang/database/factories"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserSeeder struct {
	db         *sql.DB
	repository *repositories.UserRepository
}

func NewUserSeeder(db *sql.DB) *UserSeeder {
	return &UserSeeder{
		db:         db,
		repository: repositories.NewUserRepository(db),
	}
}

func (seeder *UserSeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	users := []entities.User{}

	func() { // Seed user with specific email for testing purposes
		password, err := bcrypt.GenerateFromPassword([]byte("11112222"), bcrypt.DefaultCost)
		helper.ErrorPanic(err)

		user := entities.User{
			Id:        uuid.New(),
			Name:      "Muhammad Arfan",
			Email:     "arfan@gmail.com",
			Password:  string(password),
			CreatedAt: time.Now(),
			UpdatedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		}
		users = append(users, user)
	}()

	for i := 0; i < 50; i++ { // seed generated user by factory
		users = append(users, factories.FakeUser())
	}

	_, err := seeder.repository.Insert(ctx, users...)
	helper.ErrorPanic(err)
}
