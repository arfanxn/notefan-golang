package page_ress

import (
	"time"

	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/responses/media_ress"
	"gopkg.in/guregu/null.v4"
)

// Page resource / response
type Page struct {
	Id        string    `json:"id"`
	SpaceId   string    `json:"space_id"`
	Title     string    `json:"title"`
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt null.Time `json:"updated_at"`

	Icon  media_ress.Media `json:"icon,omitempty"`
	Cover media_ress.Media `json:"cover,omitempty"`
}

// FillFromEntity fills response from entity
func FillFromEntity(entity entities.Page) Page {
	return Page{
		Id:        entity.Id.String(),
		SpaceId:   entity.SpaceId.String(),
		Title:     entity.Title,
		Order:     entity.Order,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: null.NewTime(entity.UpdatedAt.Time, entity.UpdatedAt.Valid),
	}
}
