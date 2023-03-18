package order_bys

// Query Order Bys
const (
	Asc  = "ASC"
	Desc = "DESC"
)

// All returns slice enums
func All() []string {
	return []string{Asc, Desc}
}
