package syscall

import "fmt"

type Error struct {
	Name  string
	Errno Errno
}

func NewError(name string, errno uintptr) error {
	if errno == 0 {
		return nil
	}
	return Error{Name: name, Errno: Errno(errno)}
}

func (e Error) Error() string {
	return fmt.Sprintf("%s failed with code %d (%s)", e.Name, e.Errno, e.Errno)
}
