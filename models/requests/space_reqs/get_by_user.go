package space_reqs

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetByUser struct {
	UserId  string `json:"user_id"`
	Offset  int    `json:"offset"`
	PerPage int    `json:"page"`
	Limit   int    `json:"limit"`
	OrderBy string `json:"order_by"`
}

// Validate validates the request
func (input GetByUser) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.UserId,
			validation.Required),
		validation.Field(&input.Offset,
			validation.When(input.PerPage == 0, validation.Required)),
		validation.Field(&input.PerPage,
			validation.When(input.Offset == 0, validation.Required)),
		validation.Field(&input.Limit,
			validation.Required),
		validation.Field(&input.OrderBy, validation.Required),
	)
}
