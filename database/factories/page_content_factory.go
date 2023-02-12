package factories

import (
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func NewPageContent() entities.PageContent {
	return entities.PageContent{
		Id: uuid.New(),
		// PageId: , // Will be filled in later
		Type:      faker.Word(),
		Order:     1,
		Body:      faker.Paragraph(),
		CreatedAt: time.Now(),
		UpdatedAt: helper.RandomSQLNullTime(time.Now().AddDate(0, 0, 1)),
	}
}
