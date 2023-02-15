package factories

import (
	"strings"
	"time"

	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/entities"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func FakeNotification() entities.Notification {
	typ := strings.ReplaceAll(strings.ToUpper(faker.Word()), " ", "")

	return entities.Notification{
		Id: uuid.New(),
		// ObjectType:, // will be filled in later
		// ObjectId: , // will be filled in later
		Title:      faker.Word(),
		Type:       typ,
		Body:       faker.Paragraph(),
		ArchivedAt: helper.DBRandNullOrTime(time.Now().AddDate(0, 0, -2)),
		CreatedAt:  time.Now(),
		UpdatedAt:  helper.DBRandNullOrTime(time.Now().AddDate(0, 0, 1)),
	}
}
