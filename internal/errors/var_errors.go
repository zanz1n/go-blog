package errors

var ise = &StatusError{
	Code:     50000,
	HttpCode: 500,
	Message:  "Something went wrong",
}

var (
	ErrInternalServerError = ise
)
