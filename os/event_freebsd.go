//go:build freebsd
// +build freebsd

package os

import (
	"unsafe"

	"github.com/anton2920/gofa/bits"
	"github.com/anton2920/gofa/os/posix/freebsd"
)

type Event struct {
	Identifier  uintptr
	EventType   int16
	ActionFlags bits.Flags16
	EventNotes  bits.Flags32
	EventData   int64
	UserData    unsafe.Pointer
	_           [4]uint64
}

const (
	EventTypeRead   = freebsd.EVFILT_READ
	EventTypeWrite  = freebsd.EVFILT_WRITE
	EventTypeSignal = freebsd.EVFILT_SIGNAL
	EventTypeTimer  = freebsd.EVFILT_TIMER
)

const (
	EventFlagAdd         = bits.Flags16(freebsd.EV_ADD)
	EventFlagDelete      = bits.Flags16(freebsd.EV_DELETE)
	EventFlagOneshot     = bits.Flags16(freebsd.EV_ONESHOT)
	EventFlagTriggerEdge = bits.Flags16(freebsd.EV_CLEAR)
)

const (
	EventNoteNanoseconds  = bits.Flags32(freebsd.NOTE_NSECONDS)
	EventNoteAbsoluteTime = bits.Flags32(freebsd.NOTE_ABSTIME)
)

func (e *Event) EndOfFile() bool {
	return e.ActionFlags.Have(freebsd.EV_EOF)
}

func (e *Event) Error() freebsd.Errno {
	if e.ActionFlags.Have(freebsd.EV_ERROR) {
		return freebsd.Errno(e.EventData)
	}
	return 0
}

func CreateNewEventQueue() (Handle, error) {
	q, err := freebsd.Kqueue()
	return Handle(q), err
}

func RegisterEventsWithQueue(q Handle, chlist []Event) error {
	_, err := RegisterAndReturnPendingEventsFromQueue(q, chlist, nil, nil)
	return err
}

func ReturnPendingEventsFromQueue(q Handle, evlist []Event, t *SecondsWithNanoseconds) (int, error) {
	return RegisterAndReturnPendingEventsFromQueue(q, nil, evlist, t)
}

func RegisterAndReturnPendingEventsFromQueue(q Handle, chlist []Event, evlist []Event, t *SecondsWithNanoseconds) (int, error) {
	for i := 0; i < len(chlist); i++ {
		chlist[i].ActionFlags |= EventFlagAdd
	}
	n, err := freebsd.Kevent(int32(q), *(*[]freebsd.Kevent_t)(unsafe.Pointer(&chlist)), *(*[]freebsd.Kevent_t)(unsafe.Pointer(&evlist)), (*freebsd.Timespec)(unsafe.Pointer(t)))
	return int(n), err
}
