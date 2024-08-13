package event

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/time"
)

type platformEventQueue struct {
	kq int32

	/* Events buffer */
	events [64]Event
	head   int
	tail   int
}

type Type int16

const (
	Read   = syscall.EVFILT_READ
	Write  = syscall.EVFILT_WRITE
	Aio    = syscall.EVFILT_AIO
	Signal = syscall.EVFILT_SIGNAL
	Timer  = syscall.EVFILT_TIMER
)

type DurationUnits int

const (
	Seconds      DurationUnits = syscall.NOTE_SECONDS
	Milliseconds               = syscall.NOTE_MSECONDS
	Microseconds               = syscall.NOTE_USECONDS
	Nanoseconds                = syscall.NOTE_NSECONDS
	Absolute                   = syscall.NOTE_ABSTIME
)

type Event struct {
	Identifier uintptr
	Type       int16
	Flags      uint16
	Fflags     uint32
	Data       int
	UserData   unsafe.Pointer
	_          [4]uint
}

func (e *Event) EndOfFile() bool {
	return (e.Flags & syscall.EV_EOF) == syscall.EV_EOF
}

func (e *Event) Error() syscall.Errno {
	if (e.Flags & syscall.EV_ERROR) == syscall.EV_ERROR {
		return syscall.Errno(e.Data)
	}
	return 0
}

func platformNewEventQueue(q *Queue) error {
	kq, err := syscall.Kqueue()
	if err != nil {
		return fmt.Errorf("failed to open kernel queue: %w", err)
	}
	q.kq = kq
	return nil
}

func platformQueueAddSocket(q *Queue, l int32, request Request, trigger Trigger, userData unsafe.Pointer) error {
	var flags uint16 = syscall.EV_ADD
	if trigger == TriggerEdge {
		flags |= syscall.EV_CLEAR
	}

	if (request & RequestRead) == RequestRead {
		event := syscall.Kevent_t{Ident: uintptr(l), Filter: syscall.EVFILT_READ, Flags: flags, Udata: userData}
		if _, err := syscall.Kevent(q.kq, *(*[]syscall.Kevent_t)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&event)), Len: 1, Cap: 1})), nil, nil); err != nil {
			return fmt.Errorf("failed to request socket read event: %w", err)
		}
	}

	if (request & RequestWrite) == RequestWrite {
		event := syscall.Kevent_t{Ident: uintptr(l), Filter: syscall.EVFILT_WRITE, Flags: flags, Udata: userData}
		if _, err := syscall.Kevent(q.kq, *(*[]syscall.Kevent_t)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&event)), Len: 1, Cap: 1})), nil, nil); err != nil {
			return fmt.Errorf("failed to request socket write event: %w", err)
		}
	}

	return nil
}

func platformQueueAddSignal(q *Queue, sig syscall.Signal) error {
	event := syscall.Kevent_t{Ident: uintptr(sig), Filter: syscall.EVFILT_SIGNAL, Flags: syscall.EV_ADD}
	if _, err := syscall.Kevent(q.kq, *(*[]syscall.Kevent_t)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&event)), Len: 1, Cap: 1})), nil, nil); err != nil {
		return fmt.Errorf("failed to request signal event: %w", err)
	}
	return nil
}

func platformQueueAddTimer(q *Queue, identifier uintptr, timeout int, units DurationUnits, userData unsafe.Pointer) error {
	event := syscall.Kevent_t{Ident: identifier, Filter: syscall.EVFILT_TIMER, Flags: syscall.EV_ADD, Fflags: uint32(units), Udata: userData}
	if _, err := syscall.Kevent(q.kq, *(*[]syscall.Kevent_t)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&event)), Len: 1, Cap: 1})), nil, nil); err != nil {
		return fmt.Errorf("failed to request timer event: %w", err)
	}
	return nil
}

func platformQueueAppendEvent(q *Queue, event Event) {
	if q.head == q.tail {
		q.head = 0
		q.tail = 0
	}
	if q.tail == len(q.events)-1 {
		panic("no space left for events")
	}

	q.events[q.tail] = event
	q.tail++
}

func platformQueueClose(q *Queue) error {
	return syscall.Close(q.kq)
}

func platformQueueRequestNewEvents(q *Queue, tp *syscall.Timespec) error {
	var err error
retry:
	q.tail, err = syscall.Kevent(q.kq, nil, *(*[]syscall.Kevent_t)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&q.events[0])), Len: len(q.events), Cap: cap(q.events)})), tp)
	if err != nil {
		if err.(syscall.Error).Errno == syscall.EINTR {
			goto retry
		}
		return err
	}
	q.head = 0

	return nil
}

func platformQueueGetEvents(q *Queue, events []Event) (int, error) {
	if q.head < q.tail {
		n := copy(events, q.events[q.head:q.tail])
		q.head = 0
		q.tail = 0
		return n, nil
	}

retry:
	n, err := syscall.Kevent(q.kq, nil, *(*[]syscall.Kevent_t)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&events[0])), Len: len(events), Cap: cap(events)})), nil)
	if err != nil {
		if err.(syscall.Error).Errno == syscall.EINTR {
			goto retry
		}
	}
	return n, err
}

func platformQueueHasEvents(q *Queue) bool {
	if q.head < q.tail {
		return true
	}

	var tp syscall.Timespec
	if err := platformQueueRequestNewEvents(q, &tp); err != nil {
		return false
	}
	return q.tail > 0
}

func platformQueuePause(q *Queue, duration int64) {
	if q.head < q.tail {
		return
	}

	tp := syscall.Timespec{Sec: duration / time.NsecPerSec, Nsec: duration % time.NsecPerSec}
	platformQueueRequestNewEvents(q, &tp)
}
