package exceptions

import "errors"

var (
	// ValidationPasswordNotMatch
	ValidationPasswordNotMatch error = errors.New("Password not match")
)
