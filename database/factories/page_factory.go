package factories

import (
	"time"

	"github.com/notefan-golang/helpers/nullh"
	"github.com/notefan-golang/models/entities"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func FakePage() entities.Page {
	page := entities.Page{
		Id: uuid.New(),
		// SpaceId: ,// will be filled in later
		Title:     faker.Word(),
		Order:     1,
		CreatedAt: time.Now().AddDate(0, 0, -1),
		UpdatedAt: nullh.RandSqlNullTime(time.Now().AddDate(0, 0, 1)),
	}

	return page
}
