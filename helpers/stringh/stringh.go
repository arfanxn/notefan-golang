package stringh

import "strings"

// SliceColumnToStr format a slice of database table column names to a string
// return example "`column1`, `column2`, `column3`"
func SliceColumnToStr(columns []string) string {
	names := []string{}
	for i := 0; i < len(columns); i++ {
		names = append(names, "`"+columns[i]+"`")
	}
	return strings.Join(names, ", ")
}
