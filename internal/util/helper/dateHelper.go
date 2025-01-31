package helper

import (
	"telemetry-sale/internal/util/constant"
	"time"
)

func GetDateByPeriod(period constant.PeriodEnum) time.Time {
	var timeAgo time.Time

	switch period {
	case constant.DAY:
		now := time.Now()
		timeAgo = DateWithoutTime(now)
	case constant.WEEK:
		day := int(time.Now().Weekday()) - 1
		timeAgo = DateWithoutTime(time.Now().AddDate(0, 0, -day))
	case constant.MONTH:
		day := time.Now().Day() - 1
		timeAgo = DateWithoutTime(time.Now().AddDate(0, 0, -day))
	}

	return time.Date(timeAgo.Year(), timeAgo.Month(), timeAgo.Day(), 0, 0, 0, 0, &time.Location{})
}

func GetDateByNewPeriod(period constant.NewPeriodEnum) time.Time {
	var timeAgo time.Time

	switch period {
	case constant.Today:
		now := time.Now()
		timeAgo = DateWithoutTime(now)
	case constant.SinceMonday:
		day := int(time.Now().Weekday()) - 1
		timeAgo = DateWithoutTime(time.Now().AddDate(0, 0, -day))
	case constant.SinceFirstDayOfMonth:
		day := time.Now().Day() - 1
		timeAgo = DateWithoutTime(time.Now().AddDate(0, 0, -day))
	case constant.Last24Hours:
		timeAgo = DateWithoutTime(time.Now().AddDate(0, 0, -1))
	case constant.Last7Days:
		timeAgo = DateWithoutTime(time.Now().AddDate(0, 0, -7))
	case constant.Last30Days:
		timeAgo = DateWithoutTime(time.Now().AddDate(0, -1, 0))
	}

	return time.Date(timeAgo.Year(), timeAgo.Month(), timeAgo.Day(), 0, 0, 0, 0, &time.Location{})
}

func GetWeedDay(weekday string) string {
	var shortWeekday string

	switch weekday {
	case "Monday":
		shortWeekday = string(constant.MO)
	case "Tuesday":
		shortWeekday = string(constant.TU)
	case "Wednesday":
		shortWeekday = string(constant.WE)
	case "Thursday":
		shortWeekday = string(constant.TH)
	case "Friday":
		shortWeekday = string(constant.FR)
	case "Saturday":
		shortWeekday = string(constant.SA)
	case "Sunday":
		shortWeekday = string(constant.SU)
	}

	return shortWeekday
}

func DateWithoutTime(data time.Time) time.Time {
	return time.Date(data.Year(), data.Month(), data.Day(), 0, 0, 0, 0, data.Location())
}
