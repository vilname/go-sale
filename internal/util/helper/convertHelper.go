package helper

import (
	"database/sql"
	"strconv"
	"strings"
	"telemetry-sale/internal/util/constant"
)

func WeekDayEnumToString(weekDays []constant.Weekday) string {
	var result string
	for _, weekDay := range weekDays {
		result += string(weekDay) + ","
	}

	return strings.Trim(result, ",")
}

func Uint64ToString(ids []uint64) string {
	var result string
	for _, id := range ids {
		result += strconv.Itoa(int(id)) + ","
	}

	return strings.Trim(result, ",")
}

func Uint64ToNullInt64(changeValue uint64) sql.NullInt64 {
	var resultInt sql.NullInt64
	resultInt.Int64 = int64(changeValue)
	if changeValue > 0 {
		resultInt.Valid = true
	}

	return resultInt
}
