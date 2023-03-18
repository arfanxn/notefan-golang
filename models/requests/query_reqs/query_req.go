package query_reqs

import (
	query_order_bys "github.com/notefan-golang/enums/query/order_bys"
	"github.com/notefan-golang/helpers/sliceh"
)

// Query represents database query request
type Query struct {
	Keyword  string // saerch keyword
	Limit    int
	Offset   int64
	Withs    []string
	Wheres   map[string]any
	OrWheres map[string]any
	OrderBys map[string]string
}

func Default() Query {
	return Query{
		Offset: 0,
		Limit:  100,
		OrderBys: map[string]string{
			"created_at": query_order_bys.Desc,
		},
	}
}

// IsWith checks if the query.Withs contains the specified with
func (query *Query) IsWith(with string) bool {
	return sliceh.Contains(query.Withs, with)
}

// AddWith adds / appends to query.Withs
func (query *Query) AddWith(with string) {
	query.Withs = append(query.Withs, with)
}

// AddOrderBy adds / appends to query.OrderBys
func (query *Query) AddOrderBy(key, value string) {
	query.OrderBys[key] = value
}
