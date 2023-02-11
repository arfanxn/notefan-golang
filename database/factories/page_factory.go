package factories

import (
	"database/sql"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func NewPage(order int) entities.Page {
	page := entities.Page{
		Id:        uuid.New(),
		Title:     faker.Sentence(),
		Order:     order,
		CreatedAt: time.Now().AddDate(0, 0, -1),
	}

	updatedAt, ok := helper.Ternary(helper.BooleanRandom(), time.Now().AddDate(0, 0, -1), nil).(time.Time)
	if ok {
		page.UpdatedAt = sql.NullTime{Time: updatedAt}
	}

	return page
}
