package event

import (
	"unsafe"

	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/syscall"
)

type platformEventQueue struct {
}

/* TODO(anton2920): this is a stub, replace with actual event numbers. */
const (
	Read   = syscall.EVFILT_READ
	Write  = syscall.EVFILT_WRITE
	Aio    = syscall.EVFILT_AIO
	Signal = syscall.EVFILT_SIGNAL
	Timer  = syscall.EVFILT_TIMER
)

/* TODO(anton2920): this is a stub, replace with actual event structure. */
type Event struct {
	Identifier uintptr
	Type       int16
	Flags      uint16
	Fflags     uint32
	Data       int
	UserData   unsafe.Pointer
	_          [4]uint
}

type DurationUnits int

/* TODO(anton2920): this is a stub, replace with actual duration units. */
const (
	Seconds      = DurationUnits(syscall.NOTE_SECONDS)
	Milliseconds = syscall.NOTE_MSECONDS
	Microseconds = syscall.NOTE_USECONDS
	Nanoseconds  = syscall.NOTE_NSECONDS
	Absolute     = syscall.NOTE_ABSTIME
)

func platformNewEventQueue(q *Queue) error {
	return errors.New("not implemented")
}

func platformQueueAddSocket(q *Queue, l int32, request Request, trigger Trigger, userData unsafe.Pointer) error {
	return errors.New("not implemented")
}

func platformQueueAddSignal(q *Queue, sig syscall.Signal) error {
	return errors.New("not implemented")
}

func platformQueueAddTimer(q *Queue, identifier uintptr, timeout int, units DurationUnits, userData unsafe.Pointer) error {
	return errors.New("not implemented")
}

func platformQueueAppendEvent(q *Queue, event Event) error {
	return errors.New("not implemented")
}

func platformQueueClose(q *Queue) error {
	return errors.New("not implemented")
}

func platformQueueGetEvents(q *Queue, events []Event) (int, error) {
	return 0, errors.New("not implemented")
}

func platformQueueHasEvents(q *Queue) bool {
	return false
}

func platformQueuePause(q *Queue, duration int64) {
}
