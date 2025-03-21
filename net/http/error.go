package http

import (
	"fmt"

	"github.com/anton2920/gofa/errors"
)

type Error struct {
	StatusCode     Status
	DisplayMessage string
	LogError       error
}

var (
	NoSpaceLeft    = Error{StatusCode: StatusRequestEntityTooLarge, DisplayMessage: "no space left in the buffer", LogError: errors.New("no space left in the buffer")}
	TooManyClients = Error{StatusCode: StatusServiceUnavailable, DisplayMessage: "too many clients", LogError: errors.New("too many clients")}
)

func BadRequest(format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{StatusCode: StatusBadRequest, DisplayMessage: message, LogError: errors.WrapWithTrace(errors.New(message), 2)}
}

func Unauthorized(format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{StatusCode: StatusUnauthorized, DisplayMessage: message, LogError: errors.WrapWithTrace(errors.New(message), 2)}
}

func Forbidden(format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{StatusCode: StatusForbidden, DisplayMessage: message, LogError: errors.WrapWithTrace(errors.New(message), 2)}
}

func NotFound(format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{StatusCode: StatusNotFound, DisplayMessage: message, LogError: errors.WrapWithTrace(errors.New(message), 2)}
}

func Conflict(format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{StatusCode: StatusConflict, DisplayMessage: message, LogError: errors.WrapWithTrace(errors.New(message), 2)}
}

func RequestEntityTooLarge(format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{StatusCode: StatusRequestEntityTooLarge, DisplayMessage: message, LogError: errors.WrapWithTrace(errors.New(message), 2)}
}

func ClientError(err error) Error {
	return Error{StatusCode: StatusBadRequest, DisplayMessage: "whoops... Something went wrong. Please reload this page or try again later", LogError: errors.WrapWithTrace(err, 2)}
}

func ServerError(err error) Error {
	return Error{StatusCode: StatusInternalServerError, DisplayMessage: "whoops... Something went wrong. Please try again later", LogError: errors.WrapWithTrace(err, 2)}
}

func (err Error) Error() string {
	if err.LogError == nil {
		return "<nil>"
	}
	return err.LogError.Error()
}
