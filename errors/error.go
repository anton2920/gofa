package errors

type Error string

var ErrNotImplemented = error(New("not implemented"))

func New(msg string) Error {
	return Error(msg)
}

func (e Error) Error() string {
	return string(e)
}
