package user_reqs

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	fileReqs "github.com/notefan-golang/models/requests/file_reqs"
	fileRules "github.com/notefan-golang/rules/file_rules"
)

// UpdateProfile represents user update profile request
type UpdateProfile struct {
	Id     string        `json:"id"`
	Name   string        `json:"name"`
	Avatar fileReqs.File `json:"-"`
}

// Validate validates the Update request
func (input UpdateProfile) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Id, validation.Required),
		validation.Field(&input.Name, validation.Required, validation.Length(2, 50)),
		validation.Field(&input.Avatar,
			validation.By(fileRules.File(
				true,
				0,
				10<<20,
				[]string{"image/jpeg"}),
			),
		),
	)
}
