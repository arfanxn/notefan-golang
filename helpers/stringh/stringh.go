package stringh

import (
	"strings"
	"unicode"
)

// SliceColumnToStr format a slice of database table column names to a string
// return example "`column1`, `column2`, `column3`"
func SliceColumnToStr(columns []string) string {
	names := []string{}
	for i := 0; i < len(columns); i++ {
		names = append(names, "`"+columns[i]+"`")
	}
	return strings.Join(names, ", ")
}

func SnakeCaseToCapitalized(str string) string {
	str = strings.ReplaceAll(str, "_", " ")
	str = strings.TrimSpace(str)
	firstCharOfStr := []rune(str)[0]
	str = ReplaceAtIndex(str, 0, unicode.ToUpper(firstCharOfStr))
	return str
}

func ReplaceAtIndex(in string, i int, r rune) string {
	runes := []rune(in)
	runes[i] = r
	return string(runes)
}
