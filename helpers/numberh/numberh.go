package numberh

import (
	"math/rand"
	"time"
)

// Random generates random numbers on range
func Random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
