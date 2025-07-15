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

var (
	NoSpaceLeft    = Error{StatusCode: StatusRequestEntityTooLarge, DisplayErrorMessage: "no space left in the buffer", LogError: errors.New("no space left in the buffer")}
	TooManyClients = Error{StatusCode: StatusServiceUnavailable, DisplayErrorMessage: "too many clients", LogError: errors.New("too many clients")}
)

func (err Error) Error() string {
	if err.LogError == nil {
		return "<nil>"
	}
	return err.LogError.Error()
}

func BadRequest(format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{StatusCode: StatusBadRequest, DisplayErrorMessage: message, LogError: errors.WrapWithTrace(errors.New(message), 2)}
}

func Unauthorized(format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{StatusCode: StatusUnauthorized, DisplayErrorMessage: message, LogError: errors.WrapWithTrace(errors.New(message), 2)}
}

func Forbidden(format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{StatusCode: StatusForbidden, DisplayErrorMessage: message, LogError: errors.WrapWithTrace(errors.New(message), 2)}
}

func NotFound(format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{StatusCode: StatusNotFound, DisplayErrorMessage: message, LogError: errors.WrapWithTrace(errors.New(message), 2)}
}

func Conflict(format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{StatusCode: StatusConflict, DisplayErrorMessage: message, LogError: errors.WrapWithTrace(errors.New(message), 2)}
}

func RequestEntityTooLarge(format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{StatusCode: StatusRequestEntityTooLarge, DisplayErrorMessage: message, LogError: errors.WrapWithTrace(errors.New(message), 2)}
}

func ClientError(err error) Error {
	return Error{StatusCode: StatusBadRequest, DisplayErrorMessage: ClientDisplayErrorMessage, LogError: errors.WrapWithTrace(err, 2)}
}

func ServerError(err error) Error {
	return Error{StatusCode: StatusInternalServerError, DisplayErrorMessage: ServerDisplayErrorMessage, LogError: errors.WrapWithTrace(err, 2)}
}
