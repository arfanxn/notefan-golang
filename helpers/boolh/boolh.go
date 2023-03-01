package boolh

import (
	"math/rand"
	"time"
)

// Random returns a random boolean value
func Random() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}

func Ternary(
	condition bool,
	onTrue any,
	onFalse any,
) any {
	if condition {
		return onTrue
	} else {
		return onFalse
	}
}
