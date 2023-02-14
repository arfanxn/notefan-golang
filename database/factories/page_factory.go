package factories

import (
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func FakePage() entities.Page {
	page := entities.Page{
		Id:        uuid.New(),
		Title:     faker.Word(),
		Order:     1,
		CreatedAt: time.Now().AddDate(0, 0, -1),
		UpdatedAt: helper.DBRandNullOrTime(time.Now().AddDate(0, 0, 1)),
	}

	return page
}
