package constant

type TypeDiscountEnum string
type Weekday string

// ReasonEnum причина по кторой промо код не валиден
type ReasonEnum string

const MaxQtyLetterPromoCode = 12

const (
	PERCENT TypeDiscountEnum = "PERCENT"
	FIXED   TypeDiscountEnum = "FIXED"
	FREE    TypeDiscountEnum = "FREE"
)

const (
	MO Weekday = "MO"
	TU Weekday = "TU"
	WE Weekday = "WE"
	TH Weekday = "TH"
	FR Weekday = "FR"
	SA Weekday = "SA"
	SU Weekday = "SU"
)

const (
	NotFound       ReasonEnum = "NOT_FOUND"
	Expired        ReasonEnum = "EXPIRED"
	UseExceed      ReasonEnum = "USE_EXCEED"
	WrongTime      ReasonEnum = "WRONG_TIME"
	WrongMachine   ReasonEnum = "WRONG_MACHINE "
	InvalidProduct ReasonEnum = "INVALID_PRODUCT"
)
