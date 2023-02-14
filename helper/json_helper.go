package helper

import (
	"encoding/json"
	"io"
)

// JSONDecodeFromReader
// A helper function to decode a JSON from Reader to Go-Lang Struct
func JSONDecodeFromReader[T any](r io.Reader) (T, error) {
	decoder := json.NewDecoder(r)
	var decoded T
	err := decoder.Decode(&decoded)
	return decoded, err
}
