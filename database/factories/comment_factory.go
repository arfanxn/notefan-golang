package factories

import (
	"time"

	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/entities"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func FakeComment() entities.Comment {
	return entities.Comment{
		Id: uuid.New(),
		//CommentableType: , // will be filled in later
		//CommentableId: , // will be filled in later
		//UserId: , // will be filled in later
		Body:       faker.Sentence(),
		ResolvedAt: helper.DBRandNullOrTime(time.Now().AddDate(0, 0, -2)),
		CreatedAt:  time.Now(),
		UpdatedAt:  helper.DBRandNullOrTime(time.Now().AddDate(0, 0, 1)),
	}
}
