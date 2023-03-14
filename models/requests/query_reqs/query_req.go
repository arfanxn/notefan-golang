package query_reqs

import "github.com/notefan-golang/helpers/sliceh"

// Query represents user update profile request
type Query struct {
	Keyword  string // saerch keyword
	Limit    int
	Offset   int64
	Withs    []string
	Wheres   []map[string]any
	OrWheres []map[string]any
	OrderBys []map[string]string
}

func Default() Query {
	return Query{
		Offset: 0,
		Limit:  100,
	}
}

// IsWith checks if the query.Withs contains the specified with
func (query *Query) IsWith(with string) bool {
	return sliceh.Contains(query.Withs, with)
}
