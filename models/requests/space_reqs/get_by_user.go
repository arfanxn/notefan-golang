package space_reqs

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetByUser struct {
	UserId  string `json:"user_id"`
	Page    int64  `json:"page"`
	PerPage int    `json:"per_page"`
	OrderBy string `json:"order_by"`
}

// Validate validates the request
func (input GetByUser) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.UserId, validation.Required),
		validation.Field(&input.Page, validation.Required),
		validation.Field(&input.PerPage, validation.Required),
		validation.Field(&input.OrderBy),
	)
}
