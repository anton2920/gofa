package syscall

import "fmt"

type Error struct {
	Name  string
	Errno int
}

func NewError(name string, errno uintptr) error {
	if errno == 0 {
		return nil
	}
	return Error{Name: name, Errno: int(errno)}
}

func (e Error) Error() string {
	return fmt.Sprintf("%s failed with code %d (%s)", e.Name, e.Errno, Strerror(e.Errno))
}
