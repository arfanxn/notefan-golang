package helper

import (
	"database/sql"
	"strings"
	"time"
)

// This will return a nil nulltime or a valid nulltime
func DBRandNullOrTime(datetime time.Time) sql.NullTime {
	datetime, ok := Ternary(BooleanRandom(), datetime, nil).(time.Time)
	if !ok {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{
		Time:  datetime,
		Valid: true,
	}
}

func SliceTableColumnsToString(columns []string) string {
	for i := 0; i < len(columns); i++ {
		columns[i] = "`" + columns[i] + "`"
	}
	return strings.Join(columns, ", ")
}
