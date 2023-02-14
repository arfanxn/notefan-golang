package exceptions

import "errors"

var (
	FileNotProvided = errors.New("File not provided")
	FileInvalidType = errors.New("Invalid file type")
)
