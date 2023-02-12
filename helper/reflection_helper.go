package helper

import (
	"reflect"
)

// GetStructFieldJsonTag returns a list of struct field json tags
func GetStructFieldJsonTag(structure any) []string {
	val := reflect.ValueOf(structure)

	totalNumField := val.Type().NumField()

	tags := make([]string, totalNumField)
	for i := 0; i < totalNumField; i++ {
		tag := val.Type().Field(i).Tag.Get("json")
		tags = append(tags, tag)
	}

	return tags
}