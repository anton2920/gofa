//go:build freebsd || openbsd
// +build freebsd openbsd

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
	events []Event
	head   int
	tail   int
}

type Type int16

const (
	TypeRead   = syscall.EVFILT_READ
	TypeWrite  = syscall.EVFILT_WRITE
	TypeAio    = syscall.EVFILT_AIO
	TypeSignal = syscall.EVFILT_SIGNAL
	TypeTimer  = syscall.EVFILT_TIMER
)

type Event struct {
	Identifier uintptr
	Type       Type
	Flags      uint16
	Fflags     uint32
	Data       int64
	UserData   unsafe.Pointer
	_          [4]uint64
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
		return fmt.Errorf("failed to open kernel queue: %v", err)
	}

	q.kq = kq
	q.events = make([]Event, 64)

	return nil
}

func platformQueueAddSocket(q *Queue, s int32, request Request, trigger Trigger, userData unsafe.Pointer) error {
	var flags uint16 = syscall.EV_ADD
	if trigger == TriggerEdge {
		flags |= syscall.EV_CLEAR
	}

	events := make([]syscall.Kevent_t, 0, 2)
	if (request & RequestRead) == RequestRead {
		events = append(events, syscall.Kevent_t{Ident: uintptr(s), Filter: syscall.EVFILT_READ, Flags: flags, Udata: userData})
	}
	if (request & RequestWrite) == RequestWrite {
		events = append(events, syscall.Kevent_t{Ident: uintptr(s), Filter: syscall.EVFILT_WRITE, Flags: flags, Udata: userData})
	}

	if _, err := syscall.Kevent(q.kq, events, nil, nil); err != nil {
		return fmt.Errorf("failed to request socket events: %v", err)
	}

	return nil
}

func platformQueueAddSignal(q *Queue, sig syscall.Signal) error {
	events := make([]syscall.Kevent_t, 1)
	events[0] = syscall.Kevent_t{Ident: uintptr(sig), Filter: syscall.EVFILT_SIGNAL, Flags: syscall.EV_ADD}

	if _, err := syscall.Kevent(q.kq, events, nil, nil); err != nil {
		return fmt.Errorf("failed to request signal event: %v", err)
	}

	return nil
}

func platformQueueAddTimer(q *Queue, id uintptr, duration int64, userData unsafe.Pointer) error {
	events := make([]syscall.Kevent_t, 1)
	events[0] = syscall.Kevent_t{Ident: id, Data: duration, Filter: syscall.EVFILT_TIMER, Flags: syscall.EV_ADD, Fflags: syscall.NOTE_NSECONDS, Udata: userData}

	if _, err := syscall.Kevent(q.kq, events, nil, nil); err != nil {
		return fmt.Errorf("failed to request timer event: %v", err)
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
retry:
	n, err := syscall.Kevent(q.kq, nil, *(*[]syscall.Kevent_t)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&q.events[q.head])), Len: len(q.events) - q.head, Cap: cap(q.events) - q.head})), tp)
	if err != nil {
		if err.(syscall.Error).Errno == syscall.EINTR {
			goto retry
		}
		return err
	}

	q.tail += n
	return nil
}

func platformQueueGetEvents(q *Queue, events []Event) (int, error) {
	if q.head < q.tail {
		n := copy(events, q.events[q.head:q.tail])
		q.head += n
		if q.head >= len(q.events) {
			q.head = 0
			q.tail = 0
		}
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

	tp := syscall.Timespec{Sec: duration / time.Second, Nsec: duration % time.Second}
	platformQueueRequestNewEvents(q, &tp)
}
