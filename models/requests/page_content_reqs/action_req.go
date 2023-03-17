package page_content_reqs

import (
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Action struct {
	PageId        string `json:"page_id"`
	PageContentId string `json:"page_content_id"`
}

// Validate validates the request data
func (input Action) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.PageId, validation.Required, ozzoIs.UUID),
		validation.Field(&input.PageContentId, validation.Required, ozzoIs.UUID),
	)
}
