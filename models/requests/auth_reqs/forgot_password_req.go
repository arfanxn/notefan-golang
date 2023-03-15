package auth_reqs

import (
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ForgotPassword struct {
	Email string `json:"email"`
}

// Validate validates the Login request
func (input ForgotPassword) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Email,
			validation.Required, validation.Length(10, 50), ozzoIs.Email),
	)
}
