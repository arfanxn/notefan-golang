package factories

import (
	"database/sql"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func NewUser() entities.User {
	password, err := bcrypt.GenerateFromPassword([]byte("11112222"), bcrypt.DefaultCost)
	helper.LogFatalIfError(err)

	user := entities.User{
		Id:        uuid.New(),
		Name:      faker.Name(),
		Email:     faker.Email(),
		Password:  string(password),
		CreatedAt: time.Now().AddDate(0, 0, -1),
	}

	updatedAt, ok := helper.Ternary(helper.BooleanRandom(), time.Now().AddDate(0, 0, -1), nil).(time.Time)
	if ok {
		user.UpdatedAt = sql.NullTime{
			Time:  updatedAt,
			Valid: true,
		}
	}

	return user
}
