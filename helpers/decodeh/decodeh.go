package decodeh

import (
	"net/url"

	"github.com/gorilla/schema"
	requests "github.com/notefan-golang/models/requests/req_contracts"
)

// FormData decodes form data, return the decoded form data or panic if error happens
func FormData[T requests.ValidateableContract](formData url.Values) (input T, err error) {
	// Decode the request body
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(&input, formData)
	if err != nil {
		return
	}
	return
}
