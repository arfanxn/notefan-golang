package decodeh

import (
	"net/url"

	"github.com/gorilla/schema"
	"github.com/notefan-golang/helpers/errorh"
	requests "github.com/notefan-golang/models/requests/req_contracts"
)

// FormData decodes form data, return the decoded form data or panic if error happens
func FormData[T requests.ValidateableContract](formData url.Values) (input T) {
	// Decode the request body
	err := schema.NewDecoder().Decode(&input, formData)
	errorh.LogPanic(err)
	return input
}
