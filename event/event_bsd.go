package event

import (
	"fmt"
	"unsafe"

	"github.com/anton2920/gofa/syscall"
)

type platformEventQueue struct {
	kq int32

	/* Events buffer */
	events [1024]syscall.Kevent_t
	head   int
	tail   int
}

var keventFilter2Type = [...]Type{
	-syscall.EVFILT_READ:   Read,
	-syscall.EVFILT_WRITE:  Write,
	-syscall.EVFILT_SIGNAL: Signal,
	-syscall.EVFILT_TIMER:  Timer,
}

var eventType2Filter = [...]int16{
	Read:   syscall.EVFILT_READ,
	Write:  syscall.EVFILT_WRITE,
	Signal: syscall.EVFILT_SIGNAL,
	Timer:  syscall.EVFILT_TIMER,
}

var measurement2Note = [...]uint32{
	Seconds:      syscall.NOTE_SECONDS,
	Milliseconds: syscall.NOTE_MSECONDS,
	Microseconds: syscall.NOTE_USECONDS,
	Nanoseconds:  syscall.NOTE_NSECONDS,
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
		if _, err := syscall.Kevent(q.kq, unsafe.Slice(&event, 1), nil, nil); err != nil {
			return fmt.Errorf("failed to request socket read event: %w", err)
		}
	}

	if (request & RequestRead) == RequestRead {
		event := syscall.Kevent_t{Ident: uintptr(l), Filter: syscall.EVFILT_WRITE, Flags: flags, Udata: userData}
		if _, err := syscall.Kevent(q.kq, unsafe.Slice(&event, 1), nil, nil); err != nil {
			return fmt.Errorf("failed to request socket write event: %w", err)
		}
	}

	return nil
}

func platformQueueAddSignal(q *Queue, sig int32) error {
	event := syscall.Kevent_t{Ident: uintptr(sig), Filter: syscall.EVFILT_SIGNAL, Flags: syscall.EV_ADD}
	if _, err := syscall.Kevent(q.kq, unsafe.Slice(&event, 1), nil, nil); err != nil {
		return fmt.Errorf("failed to request signal event: %w", err)
	}
	return nil
}

func platformQueueAddTimer(q *Queue, identifier int32, timeout int, measurement DurationMeasurement, userData unsafe.Pointer) error {
	event := syscall.Kevent_t{Ident: uintptr(identifier), Filter: syscall.EVFILT_TIMER, Flags: syscall.EV_ADD, Fflags: measurement2Note[measurement], Udata: userData}
	if _, err := syscall.Kevent(q.kq, unsafe.Slice(&event, 1), nil, nil); err != nil {
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

	var flags uint16
	if event.EndOfFile {
		flags = syscall.EV_EOF
	}

	/* TODO(anton2920): this is actually incomplete for timer events. Rework that later. */
	q.events[q.tail] = syscall.Kevent_t{Ident: uintptr(event.Identifier), Filter: eventType2Filter[event.Type], Flags: flags, Data: event.Available, Udata: event.UserData}
	q.tail++
}

func platformQueueClose(q *Queue) error {
	return syscall.Close(q.kq)
}

func platformQueueRequestNewEvents(q *Queue, tp *syscall.Timespec) error {
	var err error
retry:
	q.tail, err = syscall.Kevent(q.kq, nil, unsafe.Slice(&q.events[0], len(q.events)), tp)
	if err != nil {
		if err.(syscall.Error).Errno == syscall.EINTR {
			goto retry
		}
		return err
	}
	q.head = 0

	return nil
}

func platformQueueGetEvent(q *Queue) (Event, error) {
	if q.head >= q.tail {
		if err := platformQueueRequestNewEvents(q, nil); err != nil {
			return EmptyEvent, err
		}
	}
	head := q.events[q.head]
	q.head++

	if (head.Flags & syscall.EV_ERROR) == syscall.EV_ERROR {
		return EmptyEvent, fmt.Errorf("requested event for %v failed with code %v", head.Ident, head.Data)
	}

	return Event{Type: keventFilter2Type[-head.Filter], Identifier: int32(head.Ident), Available: head.Data, UserData: head.Udata, EndOfFile: (head.Flags & syscall.EV_EOF) == syscall.EV_EOF}, nil
}

/* platformQueueGetTime returns current time in nanoseconds. */
func platformQueueGetTime() int64 {
	var tp syscall.Timespec
	syscall.ClockGettime(syscall.CLOCK_REALTIME, &tp)
	return tp.Sec*1_000_000_000 + tp.Nsec
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

	tp := syscall.Timespec{Sec: duration / 1_000_000_000, Nsec: duration % 1_000_000_000}
	platformQueueRequestNewEvents(q, &tp)
}
