package errors

var ise = &statusErrorImpl{
	code:     50000,
	httpCode: 500,
	message:  "Something went wrong",
}

var (
	ErrInternalServerError = New("internal server error")
)

var mpe = map[error]StatusError{
	ErrInternalServerError: ise,
}
