package space_reqs

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/notefan-golang/helpers/entityh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/rules/query_rules"
)

type GetByUser struct {
	UserId  string `json:"user_id"`
	Page    int64  `json:"page"`
	PerPage int    `json:"per_page"`
	OrderBy string `json:"order_by"`
	Keyword string `json:"keyword"` // the search keyword
}

// Validate validates the request
func (input GetByUser) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.UserId, validation.Required),
		validation.Field(&input.Page, validation.Required),
		validation.Field(&input.PerPage, validation.Required),
		validation.Field(&input.OrderBy,
			validation.By(query_rules.OrderBys(entityh.GetColumnNames(entities.Space{}))),
		),
		validation.Field(&input.Keyword),
	)
}
