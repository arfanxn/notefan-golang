package exceptions

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// NewValidationError create a validation error, from the given key and message
func NewValidationError(key string, message string) error {
	values := map[string]any{}
	values[key] = ""
	validationErr := validation.Validate(values, validation.Map(
		validation.Key(key, validation.Length(10, 10).Error(message)),
	),
	)
	return validationErr
}

var (
	// ValidationPasswordNotMatch
	ValidationPasswordNotMatch error = errors.New("Password not match")

	// ValidationFileSize
	ValidationFileSize error = errors.New("File size is too large or too small")

	// ValidationFileNotProvided
	ValidationFileNotProvided error = errors.New("File not provided")

	// ValidationFileMimeType
	ValidationFileMimeType error = errors.New("File mimetype is not supported")
)
