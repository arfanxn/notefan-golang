package auth_reqs

import (
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	stringRules "github.com/notefan-golang/rules/string_rules.go"
)

type Register struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// Validate validates the Register request
func (input Register) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Name,
			validation.Required, validation.Length(2, 50)),
		validation.Field(&input.Email,
			validation.Required, validation.Length(10, 50), ozzoIs.Email),
		validation.Field(&input.Password,
			validation.Required, validation.Length(8, 50),
		),
		validation.Field(&input.ConfirmPassword,
			validation.By(stringRules.PasswordMatch(input.Password)),
		),
	)
}
