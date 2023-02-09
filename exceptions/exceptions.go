package exceptions

import "errors"

var (
	DecodingError error = errors.New("Error while decoding")
	EncodingError error = errors.New("Error while encoding")
)
