package errors

import (
	"fmt"
)

type ApiError struct {
	Status Status `json:"status"`
}

func (a ApiError) Error() string {
	return fmt.Sprintf("code: %d - message: %s", a.Status.Code, a.Status.Message)
}

type RequestError struct {
	Status     Status `json:"status"`
	StatusCode int    `json:"-"`
}

func (a RequestError) Error() string {
	return fmt.Sprintf("code: %d - message: %s", a.Status.Code, a.Status.Message)
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s Status) Fmt(args ...any) Status {
	s.Message = fmt.Sprintf(s.Message, args)
	return s
}
