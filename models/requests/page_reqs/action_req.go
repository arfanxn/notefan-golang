package page_reqs

import (
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Action struct {
	SpaceId string `json:"space_id"`
	PageId  string `json:"page_id"`
}

// Validate validates the request data
func (input Action) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.SpaceId, validation.Required, ozzoIs.UUID),
		validation.Field(&input.PageId, validation.Required, ozzoIs.UUID),
	)
}
