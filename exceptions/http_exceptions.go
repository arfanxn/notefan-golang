package exceptions

import (
	"errors"
	"net/http"
)

type HTTPError struct {
	Code int
	Err  error
}

func (httpError *HTTPError) Error() string {
	return httpError.Err.Error()
}

func NewHTTPError(code int, err error) *HTTPError {
	return &HTTPError{Code: code, Err: err}
}

/**
 * 	List of HTTP exceptions
 */
var (
	HTTPNotFound = NewHTTPError(
		http.StatusNotFound, errors.New("Not found"))
	HTTPActionUnauthorized = NewHTTPError(
		http.StatusUnauthorized, errors.New("Action is unauthorized"))

	HTTPSomethingWentWrong = NewHTTPError(
		http.StatusInternalServerError, errors.New("Something went wrong"))

	HTTPValidationFailed = NewHTTPError(
		http.StatusUnprocessableEntity, errors.New("Validation failed"))

	HTTPAuthLoginFailed = NewHTTPError(
		http.StatusUnauthorized, errors.New("Email or password does not match our records"))

	HTTPAuthNotSignIn = NewHTTPError(
		http.StatusUnauthorized, errors.New("Unauthorized action, please sign in and try again"))

	HTTPAuthTokenExpired = NewHTTPError(
		http.StatusUnauthorized, errors.New("Token expired, please sign in again"))
)
