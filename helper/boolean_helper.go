package helper

import (
	"math/rand"
	"time"
)

func BoolRandom() bool {
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
