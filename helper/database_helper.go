package helper

import (
	"database/sql"
	"strings"
	"time"
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

// This will return a nil nulltime or a valid nulltime
func RandomSQLNullTime(datetime time.Time) sql.NullTime {
	datetime, ok := Ternary(BooleanRandom(), datetime, nil).(time.Time)
	if !ok {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{
		Time:  datetime,
		Valid: true,
	}
}
