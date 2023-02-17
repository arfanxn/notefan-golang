package helper

import (
	"database/sql"
	"strings"
	"time"
)

// This will return a nil nulltime or a valid nulltime
func DBRandNullOrTime(datetime time.Time) sql.NullTime {
	datetime, ok := Ternary(BoolRandom(), datetime, nil).(time.Time)
	if !ok {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{
		Time:  datetime,
		Valid: true,
	}
}

// DBSliceColumnsToStr format a slice of column names to a string
// return example "`column1`, `column2`, `column3`"
func DBSliceColumnsToStr(columns []string) string {
	names := []string{}
	for i := 0; i < len(columns); i++ {
		names = append(names, "`"+columns[i]+"`")
	}
	return strings.Join(names, ", ")
}
