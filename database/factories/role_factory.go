package factories

import (
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/notefan-golang/models/entities"
)

func FakeRole() entities.Role {
	role := entities.Role{
		Id:   uuid.New(),
		Name: faker.Word(),
	}

	return role
}
