package page_content_reqs

import (
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	page_content_types "github.com/notefan-golang/enums/page_content/types"
	"github.com/notefan-golang/helpers/sliceh"
	"github.com/notefan-golang/models/requests/file_reqs"
	"github.com/notefan-golang/rules/file_rules"
)

type Create struct {
	PageId string            `json:"space_id"`
	Type   string            `json:"type"`
	Order  int               `json:"order"`
	Body   string            `json:"body"`
	Medias []*file_reqs.File `json:"-"`
}

// Validate validates the request data
func (input Create) Validate() error {
	pageContentTypes := sliceh.Map(page_content_types.All(), func(typ string) any {
		return any(typ)
	})

	return validation.ValidateStruct(&input,
		validation.Field(&input.PageId, validation.Required, ozzoIs.UUID),
		validation.Field(&input.Type,
			validation.Required, validation.In(pageContentTypes...),
		),
		validation.Field(&input.Order, validation.Required, validation.Min(0)),
		validation.Field(&input.Body),
		validation.Field(&input.Medias, validation.Each(
			validation.By(file_rules.File(
				false,
				0,
				10<<20,
				[]string{"image/jpeg"}),
			),
		)),
	)
}
