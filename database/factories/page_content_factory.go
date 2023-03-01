package factories

import (
	"time"

	"github.com/notefan-golang/helpers/nullh"
	"github.com/notefan-golang/models/entities"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func FakePageContent() entities.PageContent {
	return entities.PageContent{
		Id: uuid.New(),
		// PageId: , // Will be filled in later
		Type:      faker.Word(),
		Order:     1,
		Body:      faker.Paragraph(),
		CreatedAt: time.Now(),
		UpdatedAt: nullh.RandSqlNullTime(time.Now().AddDate(0, 0, 1)),
	}
}
