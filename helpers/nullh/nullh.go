package nullh

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

// IntNull returns an null.Int with null value
func IntNull(datetime time.Time) null.Int {
	return null.NewInt(0, false)
}
