package space_member_reqs

import (
	validation "github.com/go-ozzo/ozzo-validation"
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
)

type Get struct {
	SpaceId  string `json:"space_id"`
	MemberId string `json:"member_id"`
	Page     int64  `json:"page"`
	PerPage  int    `json:"per_page"`
	OrderBy  string `json:"order_by"`
	Keyword  string `json:"keyword"` // the search keyword
}

func (input Get) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.SpaceId, validation.Required, ozzoIs.UUID),
		validation.Field(&input.MemberId, validation.Required, ozzoIs.UUID),
		validation.Field(&input.Page, validation.Required),
		validation.Field(&input.PerPage, validation.Required),
		validation.Field(&input.OrderBy),
		validation.Field(&input.Keyword),
	)
}
