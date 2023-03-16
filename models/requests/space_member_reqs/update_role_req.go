package space_member_reqs

import (
	validation "github.com/go-ozzo/ozzo-validation"
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
	role_names "github.com/notefan-golang/enums/role/names"
)

// UpdateRole represents Find and Delete request action
type UpdateRole struct {
	SpaceId  string `json:"id"` // the space id
	MemberId string `json:"member_id"`
	RoleName string `json:"role_name"`
}

func (input UpdateRole) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.SpaceId, validation.Required, ozzoIs.UUID),
		validation.Field(&input.MemberId, validation.Required, ozzoIs.UUID),
		validation.Field(&input.RoleName,
			validation.Required, validation.In(role_names.SpaceOwner, role_names.SpaceMember),
		),
	)
}
