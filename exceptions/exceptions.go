package exceptions

import "errors"

var (
	// Encode and Decode errors
	DecodingError error = errors.New("Error while decoding")
	EncodingError error = errors.New("Error while encoding")

	// Query errors
	DataNotFoundError error = errors.New("Data not found")

	// Validation error(s)
	ValidationError error = errors.New("Validation error")

	SomethingWentWrongError error = errors.New("Something went wrong")

	// Invalid Signing Method
	JWTInvalidSigningMethodError error = errors.New("Invalid JWT Signing Method")
)
