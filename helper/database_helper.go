package helper

import (
	"strings"
)

func BuildBulkInsertQuery(table string, totalRows int, columns ...string) string {
	query := "INSERT INTO " + table + "({COLUMN_NAMES}) VALUES "

	columnNames := []string{}
	for _, column := range columns {
		columnNames = append(columnNames, "`"+column+"`")
	}
	query = strings.Replace(query, "{COLUMN_NAMES}", strings.Join(columnNames, ", "), 1)

	valueStrs := []string{}
	for i := 0; i < int(totalRows); i++ {
		questionMarks := strings.Split(strings.Repeat("?", len(columns)), "")
		valueStrs = append(
			valueStrs,
			"("+strings.Join(questionMarks, ", ")+")",
		)
	}

	return query + strings.Join(valueStrs, ", ")
}
