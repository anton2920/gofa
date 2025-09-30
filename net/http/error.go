package http

import (
	"fmt"

	"github.com/anton2920/gofa/errors"
)

type Error struct {
	StatusCode          Status
	DisplayErrorMessage string
	LogError            error
}

var (
	ClientDisplayErrorMessage = "whoops... Something went wrong. Please reload this page or try again later"
	ServerDisplayErrorMessage = "whoops... Something went wrong. Please try again later"
)

func (err Error) Error() string {
	if err.LogError == nil {
		return "<nil>"
	}
	return err.LogError.Error()
}

func New(status Status, format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{StatusCode: status, DisplayErrorMessage: message, LogError: errors.WrapWithTrace(errors.New(message), 3)}
}

func BadRequest(format string, args ...interface{}) Error {
	return New(StatusBadRequest, format, args...)
}

func Unauthorized(format string, args ...interface{}) Error {
	return New(StatusUnauthorized, format, args...)
}

func Forbidden(format string, args ...interface{}) Error {
	return New(StatusForbidden, format, args...)
}

func NotFound(format string, args ...interface{}) Error { return New(StatusNotFound, format, args...) }

func Conflict(format string, args ...interface{}) Error { return New(StatusConflict, format, args...) }

func RequestEntityTooLarge(format string, args ...interface{}) Error {
	return New(StatusRequestEntityTooLarge, format, args...)
}

func ServiceUnavailable(format string, args ...interface{}) Error {
	return New(StatusServiceUnavailable, format, args...)
}

func ClientError(err error) Error {
	return Error{StatusCode: StatusBadRequest, DisplayErrorMessage: ClientDisplayErrorMessage, LogError: errors.WrapWithTrace(err, 2)}
}

func ServerError(err error) Error {
	return Error{StatusCode: StatusInternalServerError, DisplayErrorMessage: ServerDisplayErrorMessage, LogError: errors.WrapWithTrace(err, 2)}
}
