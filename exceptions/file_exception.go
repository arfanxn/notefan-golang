package exceptions

import "errors"

var (
	FileNotProvided = errors.New("File not provided")
	InvalidFileType = errors.New("Invalid file type")
)
