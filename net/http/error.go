package http

import (
	"fmt"

	"github.com/anton2920/gofa/errors"
)

type Error struct {
	Status
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

func NewError(status Status, format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{Status: status, DisplayErrorMessage: message, LogError: errors.WrapWithTrace(errors.Error(message), 3)}
}

func BadRequest(format string, args ...interface{}) Error {
	return NewError(StatusBadRequest, format, args...)
}

func Unauthorized(format string, args ...interface{}) Error {
	return NewError(StatusUnauthorized, format, args...)
}

func Forbidden(format string, args ...interface{}) Error {
	return NewError(StatusForbidden, format, args...)
}

func NotFound(format string, args ...interface{}) Error {
	return NewError(StatusNotFound, format, args...)
}

func MethodNotAllowed(format string, args ...interface{}) Error {
	return NewError(StatusMethodNotAllowed, format, args...)
}

func Conflict(format string, args ...interface{}) Error {
	return NewError(StatusConflict, format, args...)
}

func RequestEntityTooLarge(format string, args ...interface{}) Error {
	return NewError(StatusRequestEntityTooLarge, format, args...)
}

func InternalServerError(format string, args ...interface{}) Error {
	return NewError(StatusInternalServerError, format, args...)
}

func ServiceUnavailable(format string, args ...interface{}) Error {
	return NewError(StatusServiceUnavailable, format, args...)
}

func ClientError(err error) Error {
	return Error{Status: StatusBadRequest, DisplayErrorMessage: ClientDisplayErrorMessage, LogError: errors.WrapWithTrace(err, 2)}
}

func ServerError(err error) Error {
	return Error{Status: StatusInternalServerError, DisplayErrorMessage: ServerDisplayErrorMessage, LogError: errors.WrapWithTrace(err, 2)}
}
