package space_reqs

import (
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/notefan-golang/models/requests/file_reqs"
	"github.com/notefan-golang/rules/file_rules"
)

type Create struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Domain      string          `json:"domain"`
	Icon        *file_reqs.File `json:"-"`

	// Metadata
	UserId string `json:"user_id"`
}

// Validate validates the request data
func (input Create) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Name, validation.Required),
		validation.Field(&input.Description, validation.Required),
		validation.Field(&input.Domain, validation.Required),
		validation.Field(&input.Icon,
			validation.By(file_rules.File(
				false,
				0,
				10<<20,
				[]string{"image/jpeg"}),
			),
		),

		// Metadata validation
		validation.Field(&input.UserId, validation.Required, ozzoIs.UUID),
	)
}
