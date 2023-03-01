package nullh

import (
	"database/sql"
	"time"

	"github.com/notefan-golang/helpers/boolh"
)

// This will return a invalid nulltime or a valid nulltime
func RandSqlNullTime(datetime time.Time) sql.NullTime {
	datetime, ok := boolh.Ternary(boolh.Random(), datetime, nil).(time.Time)
	if !ok {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{
		Time:  datetime,
		Valid: true,
	}
}
