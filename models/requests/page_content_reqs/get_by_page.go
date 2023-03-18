package page_content_reqs

import (
	ozzoIs "github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/notefan-golang/helpers/entityh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/rules/query_rules"
)

type GetByPage struct {
	PageId   string   `json:"user_id"`
	Page     int64    `json:"page"`
	PerPage  int      `json:"per_page"`
	OrderBys []string `json:"order_bys"`
	Keyword  string   `json:"keyword"` // the search keyword
}

// Validate validates the request
func (input GetByPage) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.PageId, validation.Required, ozzoIs.UUID),
		validation.Field(&input.Page, validation.Required),
		validation.Field(&input.PerPage, validation.Required),
		validation.Field(&input.OrderBys,
			validation.By(query_rules.OrderBys(entityh.GetColumnNames(entities.PageContent{}))),
		),
		validation.Field(&input.Keyword),
	)
}
