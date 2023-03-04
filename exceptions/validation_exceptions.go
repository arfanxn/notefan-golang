package exceptions

import "errors"

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
