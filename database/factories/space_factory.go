package factories

import (
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func NewSpace() entities.Space {
	space := entities.Space{
		Id:          uuid.New(),
		Name:        faker.Word(),
		Description: faker.Sentence(),
		Domain:      faker.DomainName(),
		CreatedAt:   time.Now(),
		UpdatedAt:   helper.DBRandNullOrTime(time.Now().AddDate(0, 0, 1)),
	}
	return space
}
