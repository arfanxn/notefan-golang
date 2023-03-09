package query_reqs

// Query represents user update profile request
type Query struct {
	Limit    int
	Offset   int64
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
