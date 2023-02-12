package helper

import (
	"reflect"
)

// GetStructFieldJsonTag returns a list of struct field json tags
func GetStructFieldJsonTag(structure any) []string {
	val := reflect.ValueOf(structure)

	tags := []string{}
	for i := 0; i < val.Type().NumField(); i++ {
		tag := val.Type().Field(i).Tag.Get("json")
		tags = append(tags, tag)
	}

	return tags
}
