package validationh

import (
	"net/http"

	"github.com/notefan-golang/exceptions"
	reqContracts "github.com/notefan-golang/models/requests/req_contracts"
)

// ValidateStruct validates a struct and returns a HTTPError if it is not valid
func ValidateStruct[T reqContracts.ValidateableContract](structure T) error {
	// Validate the input
	validationErrs := structure.Validate()
	if validationErrs != nil {
		// Convert validation errors to http error
		err := exceptions.NewHTTPError(http.StatusUnprocessableEntity, validationErrs)
		return err
	}
	return nil
}
