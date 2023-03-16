package auth_reqs

import (
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	stringRules "github.com/notefan-golang/rules/string_rules.go"
)

type ResetPassword struct {
	Email           string `json:"email"`
	Otp             string `json:"otp"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// Validate validates the Login request
func (input ResetPassword) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Email,
			validation.Required, validation.Length(10, 50), ozzoIs.Email),
		validation.Field(&input.Otp,
			validation.Required, validation.Length(6, 6)),
		validation.Field(&input.Password,
			validation.Required, validation.Length(8, 50),
		),
		validation.Field(&input.ConfirmPassword,
			validation.By(stringRules.PasswordMatch(input.Password)),
		),
	)
}
