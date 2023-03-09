package pagination_reqs

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Pagination struct {
	Offset  int `json:"offset"`
	PerPage int `json:"page"`
	Limit   int `json:"limit"`
}

// Validate validates the request
func (input Pagination) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Offset,
			validation.When(input.PerPage == 0, validation.Required)),
		validation.Field(&input.PerPage,
			validation.When(input.Offset == 0, validation.Required)),
		validation.Field(&input.Limit,
			validation.Required),
	)
}
