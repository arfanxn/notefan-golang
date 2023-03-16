package factories

import (
	"time"

	"github.com/google/uuid"
	"github.com/notefan-golang/helpers/nullh"
	"github.com/notefan-golang/models/entities"
)

func FakeUserRoleSpace() entities.UserRoleSpace {
	urs := entities.UserRoleSpace{
		Id: uuid.New(),
		// UserId : // will be filled in later
		// RoleId : // will be filled in later
		// SpaceId : // will be filled in later
		CreatedAt: time.Now().AddDate(0, 0, -1),
		UpdatedAt: nullh.RandSqlNullTime(time.Now().AddDate(0, 0, 1)),
	}

	return urs
}
