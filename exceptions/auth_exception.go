package exceptions

import "errors"

var (
	AuthFailedToRegister = errors.New("Failed to register user")
)
