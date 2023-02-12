package helper

import (
	"database/sql"
	"time"
)

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
