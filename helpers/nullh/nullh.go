package nullh

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

// NullInt returns an null.Int with null value
func NullInt(datetime time.Time) null.Int {
	return null.NewInt(0, false)
}
