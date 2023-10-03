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
		HttpCode: 409,
		Message:  "This user already exists, maybe try a different email",
	}
	ErrAuthTokenGenFailed = &StatusError{
		Code:     50002,
		HttpCode: 500,
		Message:  "Failed to generate authentication token, try again later",
	}
	ErrInvalidAuthToken = &StatusError{
		Code:     40101,
		HttpCode: 401,
		Message:  "The authentication token is not longer valid, please login again",
	}
	ErrLoginFailed = &StatusError{
		Code: 40102,
		HttpCode: 401,
		Message: "User don't exist or passwords do not match",
	}
)
