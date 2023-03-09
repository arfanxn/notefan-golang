package stringh

import (
	"path/filepath"
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

// SliceTableColumnToStr format a slice of database table column names to a string
// return example "`table.column1`, `table.column2`, `table.column3`"
func SliceTableColumnToStr(tableName string, columns []string) string {
	names := []string{}
	for i := 0; i < len(columns); i++ {
		names = append(names, "`"+tableName+"."+columns[i]+"`")
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

func ReplaceAtIndex(str string, index int, r rune) string {
	runes := []rune(str)
	runes[index] = r
	return string(runes)
}

// FileNameWithoutExt returns the file name only (without extension and directory)
func FileNameWithoutExt(fileName string) string {
	return filepath.Base(
		fileName[:len(fileName)-len(filepath.Ext(fileName))],
	)
}
