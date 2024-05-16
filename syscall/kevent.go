package syscall

import "unsafe"

type Kevent_t struct {
	Ident  uintptr
	Filter int16
	Flags  uint16
	Fflags uint32
	Data   int
	Udata  unsafe.Pointer
	Ext    [4]uint
}

const (
	/* From <sys/event.h>. */
	EVFILT_READ   = -1
	EVFILT_WRITE  = -2
	EVFILT_VNODE  = -4 /* attached to vnodes */
	EVFILT_SIGNAL = -6 /* attached to struct proc */
	EVFILT_TIMER  = -7 /* timers */

	EV_ADD   = 0x0001 /* add event to kq (implies enable) */
	EV_CLEAR = 0x0020 /* clear event state after reporting */

	EV_ERROR = 0x4000 /* error, data contains errno */
	EV_EOF   = 0x8000 /* EOF detected */

	NOTE_WRITE = 0x0002 /* data contents changed */

	NOTE_SECONDS  = 0x00000001 /* data is seconds */
	NOTE_MSECONDS = 0x00000002 /* data is milliseconds */
	NOTE_USECONDS = 0x00000004 /* data is microseconds */
	NOTE_NSECONDS = 0x00000008 /* data is nanoseconds */
)
