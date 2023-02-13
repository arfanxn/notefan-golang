package helper

import (
	"reflect"
)

// ReflectGetStructFieldJsonTag returns a list of struct field json tags
func ReflectGetStructFieldJsonTag(structure any) []string {
	val := reflect.ValueOf(structure)

	tags := []string{}
	for i := 0; i < val.Type().NumField(); i++ {
		tag := val.Type().Field(i).Tag.Get("json")
		tags = append(tags, tag)
	}

	return tags
}

func ReflectGetTypeName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
