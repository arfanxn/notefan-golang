package space_ress

import (
	"time"

	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/responses/media_ress"
	"gopkg.in/guregu/null.v4"
)

type Space struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Domain      string    `json:"domain"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   null.Time `json:"updated_at"`

	Icon media_ress.Media `json:"icon,omitempty"`
}

func FillFromEntity(entity entities.Space) Space {
	return Space{
		Id:          entity.Id.String(),
		Name:        entity.Name,
		Description: entity.Description,
		Domain:      entity.Domain,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   null.NewTime(entity.UpdatedAt.Time, entity.UpdatedAt.Valid),
	}
}
