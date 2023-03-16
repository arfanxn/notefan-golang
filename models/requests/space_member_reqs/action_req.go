package space_member_reqs

import (
	validation "github.com/go-ozzo/ozzo-validation"
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
)

// Action represents Find and Delete request action
type Action struct {
	SpaceId  string `json:"id"` // the space id
	MemberId string `json:"member_id"`
}

func (input Action) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.SpaceId, validation.Required, ozzoIs.UUID),
		validation.Field(&input.MemberId, validation.Required, ozzoIs.UUID),
	)
}
