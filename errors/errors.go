package errors

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

type Error string

func New(msg string) Error {
	return Error(msg)
}

func (e Error) Error() string {
	return string(e)
}

type Panic struct {
	Value interface{}
	Trace []byte
}

func NewPanic(value interface{}) Panic {
	return Panic{Value: value, Trace: debug.Stack()}
}

func (e Panic) Error() string {
	buffer := make([]byte, 0, 1024)
	buffer = fmt.Appendf(buffer, "%v\n", e.Value)
	buffer = append(buffer, e.Trace...)
	return string(buffer)
}

func WrapWithTrace(err error, skip int) error {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return err
	}
	return fmt.Errorf("%s:%d: %w", file, line, err)
}
