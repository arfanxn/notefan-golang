package factories

import (
	"time"

	"github.com/notefan-golang/helpers/nullh"
	"github.com/notefan-golang/models/entities"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func FakeUser() entities.User {
	// ! Disabled due to seeding time issues
	// password, err := bcrypt.GenerateFromPassword([]byte("11112222"), bcrypt.DefaultCost)
	// errorh.LogPanic(err)

	user := entities.User{
		Id:        uuid.New(),
		Name:      faker.Name(),
		Email:     faker.Email(),
		Password:  faker.Password(),
		CreatedAt: time.Now().AddDate(0, 0, -1),
		UpdatedAt: nullh.RandSqlNullTime(time.Now().AddDate(0, 0, 1)),
	}

	return user
}
