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

// JSONStructToJSONStr parse a go-lang struct to JSON string
func JSONStructToJSONStr(data any) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
