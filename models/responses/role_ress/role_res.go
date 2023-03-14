package role_ress

import (
	"github.com/notefan-golang/models/entities"
)

type Role struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func FillFromEntity(entity entities.Role) Role {
	return Role{
		Id:   entity.Id.String(),
		Name: entity.Name,
	}
}
