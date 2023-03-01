package factories

import (
	"time"

	"github.com/notefan-golang/helpers/nullh"
	"github.com/notefan-golang/models/entities"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func FakeSpace() entities.Space {
	space := entities.Space{
		Id:          uuid.New(),
		Name:        faker.Word(),
		Description: faker.Sentence(),
		Domain:      faker.DomainName(),
		CreatedAt:   time.Now(),
		UpdatedAt:   nullh.RandSqlNullTime(time.Now().AddDate(0, 0, 1)),
	}
	return space
}
