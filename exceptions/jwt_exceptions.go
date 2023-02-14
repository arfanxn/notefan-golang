package exceptions

import "errors"

var (
	// Invalid Signing Method
	JWTInvalidSigningMethod error = errors.New("Invalid JWT Signing Method")
)
