package seeders

import (
	"database/sql"
	"notefan-golang/database/factories"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"runtime"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserSeeder struct {
	db        *sql.DB
	tableName string
	repo      *repositories.UserRepo
}

func NewUserSeeder(db *sql.DB) *UserSeeder {
	return &UserSeeder{
		db:        db,
		tableName: "users",
		repo:      repositories.NewUserRepo(db),
	}
}

func (seeder *UserSeeder) Run() {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunning(pc)
	defer printFinishRunning(pc)

	// ---- Begin ----
	totalRows := 50
	valueArgs := []any{}

	for i := 0; i < totalRows; i++ {
		if i == 0 {
			password, err := bcrypt.GenerateFromPassword([]byte("11112222"), bcrypt.DefaultCost)
			helper.PanicIfError(err)

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
			valueArgs = append(valueArgs,
				user.Id.String(), user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt,
			)

			continue
		}

		user := factories.NewUser()
		valueArgs = append(
			valueArgs,
			user.Id.String(), user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	}

	query := helper.BuildBulkInsertQuery(seeder.tableName, totalRows,
		`id`, `name`, `email`, `password`, `created_at`, `updated_at`)

	stmt, err := seeder.db.Prepare(query)
	helper.LogFatalIfError(err)

	_, err = stmt.Exec(valueArgs...)
	helper.LogFatalIfError(err)
}
