package page_content_reqs

import (
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/notefan-golang/models/requests/file_reqs"
	"github.com/notefan-golang/rules/file_rules"
)

type Create struct {
	PageId string            `json:"space_id"`
	Type   string            `json:"type"`
	Order  int               `json:"order"`
	Body   string            `json:"body"`
	Medias []*file_reqs.File `json:"-"`
}

// Validate validates the request data
func (input Create) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.PageId, validation.Required, ozzoIs.UUID),
		validation.Field(&input.Type, validation.Required, validation.Length(1, 50)),
		validation.Field(&input.Order, validation.Required, validation.Min(0)),
		validation.Field(&input.Body),
		validation.Field(&input.Medias, validation.Each(
			validation.By(file_rules.File(
				false,
				0,
				10<<20,
				[]string{"image/jpeg"}),
			),
		)),
	)
}
