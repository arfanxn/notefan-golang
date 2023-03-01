package jsonh

import (
	"encoding/json"
	"io"
)

// DecodeFromReader
// A helper function to decode a JSON from Reader to Go-Lang Struct
func DecodeFromReader[T any](r io.Reader) (T, error) {
	decoder := json.NewDecoder(r)
	var decoded T
	err := decoder.Decode(&decoded)
	return decoded, err
}

// ToJsonStr parse any to JSON string
func ToJsonStr(data any) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
