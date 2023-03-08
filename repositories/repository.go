package repositories

import (
	"bytes"
	"strings"
)

// buildBatchInsertQuery builds a batch insert query
func buildBatchInsertQuery(tableName string, totalRows int, columnNames ...string) string {
	replaceable := "{COLUMN_NAMES}"
	query := "INSERT INTO " + tableName + "(" + replaceable + ") VALUES "

	formattedColumnNames := []string{}
	for _, columnName := range columnNames {
		formattedColumnNames = append(formattedColumnNames, "`"+columnName+"`")
	}
	query = strings.Replace(query, replaceable, strings.Join(formattedColumnNames, ", "), 1)

	valueStrs := []string{}
	for i := 0; i < int(totalRows); i++ {
		questionMarks := strings.Split(strings.Repeat("?", len(columnNames)), "")
		valueStrs = append(
			valueStrs,
			"("+strings.Join(questionMarks, ", ")+")",
		)
	}

	return query + strings.Join(valueStrs, ", ")
}

// buildUpdateQuery builds an update query
func buildUpdateQuery(tableName string, columnNames ...string) string {
	queryBuf := bytes.NewBufferString("UPDATE ")
	queryBuf.WriteString(tableName)
	queryBuf.WriteString(" SET ")

	formattedSetables := []string{}
	for _, columnName := range columnNames {
		formattedSetables = append(formattedSetables, "`"+columnName+"` = ?")
	}
	queryBuf.WriteString(strings.Join(formattedSetables, ", "))

	return queryBuf.String()
}

// buildBatchDeleteQueryByIds builds a batch delete query by given ids
func buildBatchDeleteQueryByIds(tableName string, ids ...string) (string, []any) {
	var valueArgs []any
	queryBuf := bytes.NewBufferString("DELETE FROM ")
	queryBuf.WriteString(tableName)
	queryBuf.WriteString(" WHERE ")
	for index, id := range ids {
		valueArgs = append(valueArgs, id)
		queryBuf.WriteString("`id` = ?")
		if (index + 1) != len(ids) { // if current is not last index
			queryBuf.WriteString(" OR ")
		}
	}
	return queryBuf.String(), valueArgs
}
