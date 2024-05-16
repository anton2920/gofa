package syscall

import "fmt"

type Error struct {
	Name  string
	Errno int
}

const (
	/* From <errno.h>. */
	ENOENT      = 2      /* No such file or directory */
	EINTR       = 4      /* Interrupted system call */
	EEXIST      = 17     /* File exists */
	EPIPE       = 32     /* Broken pipe */
	EAGAIN      = 35     /* Resource temporarily unavailable */
	EWOULDBLOCK = EAGAIN /* Operation would block */
	EINPROGRESS = 36     /* Operation now in progress */
	EOPNOTSUPP  = 45     /* Operation not supported */
	ECONNRESET  = 54     /* Connection reset by peer */
	ENOSYS      = 78     /* Function not implemented */
)

func NewError(name string, errno uintptr) error {
	if errno == 0 {
		return nil
	}
	return Error{Name: name, Errno: int(errno)}
}

/* TODO(anton2920): add strerror(). */
func (e Error) Error() string {
	return fmt.Sprintf("%s failed with code %d", e.Name, e.Errno)
}
