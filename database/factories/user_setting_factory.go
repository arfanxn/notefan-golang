package factories

import (
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func NewUserSetting() entities.UserSetting {
	return entities.UserSetting{
		Id: uuid.New(),
		// UserId  // will be fill manually
		Key:       faker.Word(),
		Value:     faker.Word(),
		CreatedAt: time.Now(),
		UpdatedAt: helper.RandomSQLNullTime(time.Now().AddDate(0, 0, 1)),
	}
}
