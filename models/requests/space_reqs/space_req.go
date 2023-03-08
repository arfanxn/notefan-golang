package space_reqs

import (
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/notefan-golang/models/requests/file_reqs"
	"github.com/notefan-golang/rules/file_rules"
)

type Space struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Domain      string          `json:"domain"`
	Icon        *file_reqs.File `json:"-"`
}

// Validate validates the request data
func (input Space) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Id, validation.Required, ozzoIs.UUID),
		validation.Field(&input.Name),
		validation.Field(&input.Description),
		validation.Field(&input.Domain),
		validation.Field(&input.Icon,
			validation.By(file_rules.File(
				false,
				0,
				10<<20,
				[]string{"image/jpeg"}),
			),
		),
	)
}
