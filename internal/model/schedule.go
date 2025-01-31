package model

import (
	"database/sql"
	"telemetry-sale/internal/util/constant"
)

type NewSchedule struct {
	Name    string             `json:"name" validate:"required"`
	From    string             `json:"from" example:"10:00"`
	To      string             `json:"to" example:"15:00"`
	Weekday []constant.Weekday `json:"weekday"`
}

type Schedule struct {
	Id      uint64             `json:"id"`
	Name    string             `json:"name" validate:"required"`
	From    string             `json:"from"`
	To      string             `json:"to"`
	Weekday []constant.Weekday `json:"weekday"`
}

type ScheduleResult struct {
	Id      sql.NullInt64      `json:"id"`
	Name    sql.NullString     `json:"name" validate:"required"`
	From    sql.NullString     `json:"from"`
	To      sql.NullString     `json:"to"`
	Weekday []constant.Weekday `json:"weekday"`
}

type ScheduleCheckPromoCode struct {
	From    *string `json:"from"`
	To      *string `json:"to"`
	Weekday *string `json:"weekday"`
}
