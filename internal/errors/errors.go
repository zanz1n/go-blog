package errors

type StatusError struct {
	Code     uint
	HttpCode int
	Message  string
}

func (s *StatusError) Error() string {
	return s.Message
}

type ErrorBody struct {
	Message   string `json:"message"`
	ErrorCode uint   `json:"error_code"`
}
