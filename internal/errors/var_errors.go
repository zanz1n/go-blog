package errors

var ise = &StatusError{
	Code:     50000,
	HttpCode: 500,
	Message:  "Something went wrong",
}

var (
	ErrInternalServerError = ise
	ErrUserNotFound        = &StatusError{
		Code:     40401,
		HttpCode: 404,
		Message:  "The user could not be found",
	}
	ErrUserFetchFailed = &StatusError{
		Code:     50001,
		HttpCode: 500,
		Message:  "Something went wrong while fetching user data, try again later",
	}
	ErrUserAlreadyExists = &StatusError{
		Code:     40901,
		HttpCode: 4009,
		Message:  "This user already exists, maybe try a different email",
	}
)
