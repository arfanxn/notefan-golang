package page_content_ress

import (
	"time"

	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/responses/media_ress"
	"gopkg.in/guregu/null.v4"
)

// PageContent resource / response
type PageContent struct {
	Id        string    `json:"id"`
	PageId    string    `json:"page_id"`
	Type      string    `json:"type"`
	Order     int       `json:"order"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt null.Time `json:"updated_at"`

	Medias []media_ress.Media `json:"medias,omitempty"`
}

// FillFromEntity fills response from entity
func FillFromEntity(entity entities.PageContent) PageContent {
	return PageContent{
		Id:        entity.Id.String(),
		PageId:    entity.PageId.String(),
		Type:      entity.Type,
		Order:     entity.Order,
		Body:      entity.Body,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: null.NewTime(entity.UpdatedAt.Time, entity.UpdatedAt.Valid),
	}
}
