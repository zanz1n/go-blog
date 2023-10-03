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
	Message   string `json:"Message"`
	ErrorCode uint   `json:"error_code"`
}
