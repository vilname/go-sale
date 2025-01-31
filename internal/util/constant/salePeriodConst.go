package constant

type PeriodEnum string
type NewPeriodEnum string

const (
	DAY   PeriodEnum = "DAY"
	WEEK  PeriodEnum = "WEEK"
	MONTH PeriodEnum = "MONTH"
)

const (
	Last24Hours          NewPeriodEnum = "Last24Hours"
	Last7Days            NewPeriodEnum = "Last7Days"
	Last30Days           NewPeriodEnum = "Last30Days"
	SinceMonday          NewPeriodEnum = "SinceMonday"
	SinceFirstDayOfMonth NewPeriodEnum = "SinceFirstDayOfMonth"
	Today                NewPeriodEnum = "Today"
)
