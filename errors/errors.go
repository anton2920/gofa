package errors

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

type Error string

var (
	NotImplemented = New("not implemented")
)

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
	return fmt.Sprintf("%v\n%s", e.Value, e.Trace)
}

func WrapWithTrace(err error, skip int) error {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return err
	}
	return fmt.Errorf("%s:%d: %v", file, line, err)
}
