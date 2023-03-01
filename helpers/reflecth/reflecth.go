package reflecth

import "reflect"

// GetFieldTag returns a list of struct field tags
func GetFieldTag(structure any, tagname string) []string {
	val := reflect.ValueOf(structure)

	tags := []string{}
	for i := 0; i < val.Type().NumField(); i++ {
		tag := val.Type().Field(i).Tag.Get(tagname)
		tags = append(tags, tag)
	}

	return tags
}

// GetFieldJsonTag returns a list of struct field json tags
func GetFieldJsonTag(structure any) []string {
	return GetFieldTag(structure, "json")
}

// GetTypeName get name of type of the variable
func GetTypeName(myvar any) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
