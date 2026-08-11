package errors

type WrappedError struct {
	Message    string
	InnerError error
}

func Wrap(msg string, ierr error) WrappedError {
	return WrappedError{msg, ierr}
}

/* This function is just for compatibility with 'error' interface. */
func (we WrappedError) Error() string {
	panic("you are supposed to call 'errors.Dump', not this")
	return ""
}
