package errors

import (
	"github.com/goccy/go-json"
)

func GetErrorMessage(code uint) (e StatusError, ok bool) {
	ok = false

	for _, e = range mpe {
		if e.CustomCode() == code {
			return e, true
		}
	}

	return e, false
}

func GetStatusErr(key error) StatusError {
	v, ok := mpe[key]

	if !ok {
		return &statusErrorImpl{
			code:     50000,
			httpCode: 500,
			message:  "Unknown err: " + key.Error(),
		}
	}

	return v
}

type statusErrorImpl struct {
	code     uint
	httpCode int
	message  string
}

func (e *statusErrorImpl) Message() string {
	return e.message
}

func (e *statusErrorImpl) CustomCode() uint {
	return e.code
}

func (e *statusErrorImpl) HttpCode() int {
	return e.httpCode
}

type StatusError interface {
	Message() string
	CustomCode() uint
	HttpCode() int
}

type errorImpl struct {
	m string
}

func (e *errorImpl) Error() string {
	return e.m
}

func New(text string) error {
	return &errorImpl{m: text}
}

type ErrorBody struct {
	Message   string `json:"message"`
	ErrorCode uint   `json:"error_code"`
}

func (e *ErrorBody) Marshal() []byte {
	buf, err := json.Marshal(e)

	if err != nil {
		return []byte("{\"message\":\"Failed to marshal response body\",\"error_code\":5000}")
	}

	return buf
}
