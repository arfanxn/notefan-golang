package factories

import (
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func NewUser() entities.User {
	// ! Disabled due to seeding time issues
	// password, err := bcrypt.GenerateFromPassword([]byte("11112222"), bcrypt.DefaultCost)
	// helper.LogFatalIfError(err)

	user := entities.User{
		Id:        uuid.New(),
		Name:      faker.Name(),
		Email:     faker.Email(),
		Password:  faker.Password(),
		CreatedAt: time.Now().AddDate(0, 0, -1),
		UpdatedAt: helper.RandomSQLNullTime(time.Now().AddDate(0, 0, 1)),
	}

	return user
}
