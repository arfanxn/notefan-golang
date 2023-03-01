package factories

import (
	"time"

	"github.com/notefan-golang/helpers/nullh"
	"github.com/notefan-golang/models/entities"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func FakeUserSetting() entities.UserSetting {
	return entities.UserSetting{
		Id: uuid.New(),
		// UserId  // will be fill manually
		Key:       faker.Word(),
		Value:     faker.Word(),
		CreatedAt: time.Now(),
		UpdatedAt: nullh.RandSqlNullTime(time.Now().AddDate(0, 0, 1)),
	}
}
