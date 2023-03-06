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
