package auth_reqs

import (
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates the Login request
func (input Login) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Email,
			validation.Required, validation.Length(10, 50), ozzoIs.Email),
		validation.Field(&input.Password,
			validation.Required, validation.Length(8, 50)),
	)
}
