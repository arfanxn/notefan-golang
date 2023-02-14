package factories

import (
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func FakeUser() entities.User {
	// ! Disabled due to seeding time issues
	// password, err := bcrypt.GenerateFromPassword([]byte("11112222"), bcrypt.DefaultCost)
	// helper.ErrorLogFatal(err)

	user := entities.User{
		Id:        uuid.New(),
		Name:      faker.Name(),
		Email:     faker.Email(),
		Password:  faker.Password(),
		CreatedAt: time.Now().AddDate(0, 0, -1),
		UpdatedAt: helper.DBRandNullOrTime(time.Now().AddDate(0, 0, 1)),
	}

	return user
}
