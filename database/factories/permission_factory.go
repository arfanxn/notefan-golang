package factories

import (
	"github.com/notefan-golang/models/entities"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func FakePermission() entities.Permission {
	perm := entities.Permission{
		Id:   uuid.New(),
		Name: faker.Name(),
	}

	return perm
}
