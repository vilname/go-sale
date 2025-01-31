package constant

type ErrorType string
type ErrorMessage string

const (
	MaxAttemptGenerateCode ErrorMessage = "The number of attempts exceeded"
)

const (
	MaxAttempt ErrorType = "MAX_ATTEMPT"
	DateError  ErrorType = "DATE_ERROR"
)
