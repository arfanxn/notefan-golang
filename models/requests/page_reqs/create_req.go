package page_reqs

import (
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/notefan-golang/models/requests/file_reqs"
	"github.com/notefan-golang/rules/file_rules"
)

type Create struct {
	SpaceId string          `json:"space_id"`
	Title   string          `json:"title"`
	Order   int             `json:"order"`
	Icon    *file_reqs.File `json:"-"`
	Cover   *file_reqs.File `json:"-"`
}

// Validate validates the request data
func (input Create) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.SpaceId, validation.Required, ozzoIs.UUID),
		validation.Field(&input.Title, validation.Required, validation.Length(1, 50)),
		validation.Field(&input.Order, validation.Required, validation.Min(0)),
		validation.Field(&input.Icon,
			validation.By(file_rules.File(
				false,
				0,
				10<<20,
				[]string{"image/jpeg"}),
			),
		),
		validation.Field(&input.Cover,
			validation.By(file_rules.File(
				false,
				0,
				10<<20,
				[]string{"image/jpeg"}),
			),
		),
	)
}
