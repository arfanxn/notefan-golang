package pagination_ress

import "gopkg.in/guregu/null.v4"

type Pagination[T any] struct {
	// The total of rows in database table that match the query string
	Total int64 `json:"total"`

	// Page information
	PerPage     int   `json:"per_page"`
	CurrentPage int64 `json:"current_page"`
	LastPage    int64 `json:"last_page"`

	// Pagination url navigations
	FirstPageUrl string      `json:"first_page_url"`
	PrevPageUrl  null.String `json:"prev_page_url"`
	NextPageUrl  null.String `json:"next_page_url"`
	LastPageUrl  string      `json:"last_page_url"`
	Items        []T         `json:"items"`
}
