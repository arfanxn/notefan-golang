package factories

import (
	"time"

	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/entities"

	"github.com/google/uuid"
)

func FakeCommentReaction() entities.CommentReaction {
	return entities.CommentReaction{
		Id: uuid.New(),
		//CommentId: , // will be filled in later
		//UserId: , // will be filled in later
		CreatedAt: time.Now(),
		UpdatedAt: helper.DBRandNullOrTime(time.Now().AddDate(0, 0, 1)),
	}
}
