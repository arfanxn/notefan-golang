package space_member_reqs

import (
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Invite struct {
	SpaceId string `json:"space_id"`
	Email   string `json:"email"`
}

func (input Invite) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.SpaceId, validation.Required, ozzoIs.UUID),
		validation.Field(&input.Email, validation.Required, ozzoIs.Email),
	)
}
