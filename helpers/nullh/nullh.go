package nullh

import (
	"gopkg.in/guregu/null.v4"
)

// NullInt returns an null.Int with null value
func NullInt() null.Int {
	return null.NewInt(0, false)
}
