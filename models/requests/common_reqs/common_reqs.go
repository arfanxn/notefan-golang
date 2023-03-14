package common_reqs

import (
	validation "github.com/go-ozzo/ozzo-validation"
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
)

type UUID struct {
	Id string `json:"id"`
}

func (input UUID) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Id, validation.Required, ozzoIs.UUID),
	)
}

type UUIDPagination struct {
	Id      string `json:"id"`
	Page    int64  `json:"page"`
	PerPage int    `json:"per_page"`
	OrderBy string `json:"order_by"`
	Keyword string `json:"keyword"` // the search keyword
}

func (input UUIDPagination) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Id, validation.Required, ozzoIs.UUID),
		validation.Field(&input.Page, validation.Required),
		validation.Field(&input.PerPage, validation.Required),
		validation.Field(&input.OrderBy),
		validation.Field(&input.Keyword),
	)
}
