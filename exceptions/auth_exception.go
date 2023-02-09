package exceptions

import "errors"

var (
	AuthFailedRegister = errors.New("Failed to register user")
	AuthFailedLogin    = errors.New("Email or password does not match our records")
)
