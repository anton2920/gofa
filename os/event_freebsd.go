//go:build freebsd
// +build freebsd

package os

import (
	"unsafe"

	"github.com/anton2920/gofa/bits"
	"github.com/anton2920/gofa/os/posix/freebsd"
)

type EventType int16

const (
	EventTypeRead   = EventType(freebsd.EVFILT_READ)
	EventTypeWrite  = EventType(freebsd.EVFILT_WRITE)
	EventTypeSignal = EventType(freebsd.EVFILT_SIGNAL)
	EventTypeTimer  = EventType(freebsd.EVFILT_TIMER)
)

type Event struct {
	Identifier  uintptr
	EventType   EventType
	ActionFlags bits.Flags16
	EventFlags  bits.Flags32
	EventData   int64
	UserData    unsafe.Pointer
	_           [4]uint64
}

const (
	EventQueueActionAdd                      = bits.Flags16(freebsd.EV_ADD)
	EventQueueActionDelete                   = bits.Flags16(freebsd.EV_DELETE)
	EventQueueActionShotOnce                 = bits.Flags16(freebsd.EV_ONESHOT)
	EventQueueActionResetStateAfterRetrieval = bits.Flags16(freebsd.EV_CLEAR)
)

const (
	EventNoteSeconds      = bits.Flags32(freebsd.NOTE_SECONDS)
	EventNoteMilliseconds = bits.Flags32(freebsd.NOTE_MSECONDS)
	EventNoteMicroseconds = bits.Flags32(freebsd.NOTE_USECONDS)
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

func RegisterAndReturnPendingEventsFromQueue(q Handle, chlist []Event, evlist []Event, t *SecondsWithNanoseconds) (int, error) {
	for i := 0; i < len(chlist); i++ {
		chlist[i].ActionFlags |= EventQueueActionAdd
	}
	n, err := freebsd.Kevent(int32(q), *(*[]freebsd.Kevent_t)(unsafe.Pointer(&chlist)), *(*[]freebsd.Kevent_t)(unsafe.Pointer(&evlist)), (*freebsd.Timespec)(unsafe.Pointer(t)))
	return int(n), err
}
