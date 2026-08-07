package event_

import (
	"unsafe"

	"github.com/anton2920/gofa/bits"
	"github.com/anton2920/gofa/cpu"
	"github.com/anton2920/gofa/event"
	"github.com/anton2920/gofa/os"
	"github.com/anton2920/gofa/time"
)

type Queue struct {
	KernelQueue os.Handle
	LastSync    cpu.Cycles

	events [64]os.Event
	head   int
	tail   int
}

var units2flags = map[int]bits.Flags32{
	time.Second:      os.EventNoteSeconds,
	time.Millisecond: os.EventNoteMilliseconds,
	time.Microsecond: os.EventNoteMicroseconds,
	time.Nanosecond:  os.EventNoteNanoseconds,
}

func (q *Queue) Init() error {
	kq, err := os.CreateNewEventQueue()
	if err != nil {
		return err
	}

	q.KernelQueue = kq
	return nil
}

func (q *Queue) AddFile(f os.Handle, request bits.Flags, trigger int, userData unsafe.Pointer) error {
	var flags bits.Flags16
	if trigger == event.TriggerEdge {
		flags |= os.EventQueueActionResetStateAfterRetrieval
	}

	events := make([]os.Event, 0, 2)
	if request.Has(event.RequestRead) {
		events = append(events, os.Event{Identifier: uintptr(f), EventType: os.EventTypeRead, ActionFlags: flags, UserData: userData})
	}
	if request.Has(event.RequestWrite) {
		events = append(events, os.Event{Identifier: uintptr(f), EventType: os.EventTypeWrite, ActionFlags: flags, UserData: userData})
	}

	return os.RegisterEventsWithQueue(q.KernelQueue, events)
}

func (q *Queue) AddSignals(sigs ...os.Signal) error {
	events := make([]os.Event, 0, 64)

	for len(sigs) > 0 {
		events = events[:0]

		var i int
		for i < len(sigs) {
			events = append(events, os.Event{Identifier: uintptr(sigs[i]), EventType: os.EventTypeSignal})
			i++
		}

		if err := os.RegisterEventsWithQueue(q.KernelQueue, events); err != nil {
			return err
		}

		sigs = sigs[i:]
	}

	return nil
}

func (q *Queue) AddAndIgnoreSignals(sigs ...os.Signal) error {
	if err := q.AddSignals(sigs...); err != nil {
		return err
	}
	for i := 0; i < len(sigs); i++ {
		if err := os.InstallSignalHandler(sigs[i], os.IgnoreSignalHandler); err != nil {
			return err
		}
	}
	return nil
}

func (q *Queue) AddAndIgnoreTerminateSignals() error {
	return q.AddAndIgnoreSignals(os.SignalHangup, os.SignalInterrupt, os.SignalTerminate)
}

func (q *Queue) AddPeriodicTimer(id uintptr, quantity int, units int, userData unsafe.Pointer) error {
	events := make([]os.Event, 1)
	events[0] = os.Event{Identifier: id, EventData: int64(quantity), EventType: os.EventTypeTimer, EventFlags: units2flags[units], UserData: userData}
	return os.RegisterEventsWithQueue(q.KernelQueue, events)
}

func (q *Queue) AddTimerAt(id uintptr, at int, units int, userData unsafe.Pointer) error {
	events := make([]os.Event, 1)
	events[0] = os.Event{Identifier: id, EventData: int64(at), EventType: os.EventTypeTimer, ActionFlags: os.EventQueueActionResetStateAfterRetrieval, EventFlags: units2flags[units] | os.EventNoteAbsoluteTime, UserData: userData}
	return os.RegisterEventsWithQueue(q.KernelQueue, events)
}

func (q *Queue) Close() error {
	return os.CloseHandle(q.KernelQueue)
}

func (q *Queue) GetEvents(events []os.Event) (int, error) {
	if q.head < q.tail {
		n := copy(events, q.events[q.head:q.tail])
		q.head += n
		if q.head == q.tail {
			q.head = 0
			q.tail = 0
		}
		return n, nil
	}
	return os.ReturnPendingEventsFromQueue(q.KernelQueue, events, nil)
}

func (q *Queue) requestNewEvents(t *os.SecondsWithNanoseconds) error {
	n, err := os.ReturnPendingEventsFromQueue(q.KernelQueue, q.events[:], t)
	if err != nil {
		return err
	}

	q.tail += n
	return nil
}

func (q *Queue) HasEvents() bool {
	if q.head < q.tail {
		return true
	}

	var t os.SecondsWithNanoseconds
	if err := q.requestNewEvents(&t); err != nil {
		return false
	}

	return q.tail > 0
}

func (q *Queue) SyncFPS(fps int) {
	now := cpu.ReadPerformanceCounter()
	durationBetweenPauses := now - q.LastSync
	targetRate := int64(float64(time.Second/time.Millisecond) / float64(fps) * float64(time.Millisecond))

	duration := targetRate - durationBetweenPauses.ToNanosecondsTruncated()
	if duration > 0 {
		if q.head >= q.tail {
			t := os.SecondsWithNanoseconds{Seconds: int(duration / time.Second), Nanoseconds: int(duration % time.Second)}
			q.requestNewEvents(&t)
		}
		now = cpu.ReadPerformanceCounter()
	}
	q.LastSync = now
}
