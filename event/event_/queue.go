package event_

import (
	"unsafe"

	"github.com/anton2920/gofa/bits"
	"github.com/anton2920/gofa/context"
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

func (q *Queue) Init(ctx *context.Context) bool {
	kq, ok := os.CreateNewEventQueue(ctx)
	if !ok {
		return false
	}

	q.KernelQueue = kq
	return true
}

func (q *Queue) AddFile(ctx *context.Context, f os.Handle, request bits.Flags, trigger int, userData unsafe.Pointer) bool {
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

	return os.RegisterEventsWithQueue(ctx, q.KernelQueue, events)
}

func (q *Queue) AddSignals(ctx *context.Context, sigs ...os.Signal) bool {
	events := make([]os.Event, 0, 64)

	for len(sigs) > 0 {
		events = events[:0]

		var i int
		for i < len(sigs) {
			events = append(events, os.Event{Identifier: uintptr(sigs[i]), EventType: os.EventTypeSignal})
			i++
		}

		if !os.RegisterEventsWithQueue(ctx, q.KernelQueue, events) {
			return false
		}

		sigs = sigs[i:]
	}

	return true
}

func (q *Queue) AddAndIgnoreSignals(ctx *context.Context, sigs ...os.Signal) bool {
	if !q.AddSignals(ctx, sigs...) {
		return false
	}

	for i := 0; i < len(sigs); i++ {
		if !os.InstallSignalHandler(ctx, sigs[i], os.IgnoreSignalHandler) {
			return false
		}
	}

	return true
}

func (q *Queue) AddAndIgnoreTerminationSignals(ctx *context.Context) bool {
	return q.AddAndIgnoreSignals(ctx, os.SignalHangup, os.SignalInterrupt, os.SignalTerminate)
}

func units2flags(units int) bits.Flags32 {
	switch units {
	case time.Second:
		return os.EventNoteSeconds
	case time.Millisecond:
		return os.EventNoteMilliseconds
	case time.Microsecond:
		return os.EventNoteMicroseconds
	case time.Nanosecond:
		return os.EventNoteNanoseconds
	}
	return 0
}

func (q *Queue) AddPeriodicTimer(ctx *context.Context, id uintptr, quantity int, units int, userData unsafe.Pointer) bool {
	events := make([]os.Event, 1)
	events[0] = os.Event{Identifier: id, EventData: int64(quantity), EventType: os.EventTypeTimer, EventFlags: units2flags(units), UserData: userData}
	return os.RegisterEventsWithQueue(ctx, q.KernelQueue, events)
}

func (q *Queue) AddTimerAt(ctx *context.Context, id uintptr, at int, units int, userData unsafe.Pointer) bool {
	events := make([]os.Event, 1)
	events[0] = os.Event{Identifier: id, EventData: int64(at), EventType: os.EventTypeTimer, ActionFlags: os.EventQueueActionResetStateAfterRetrieval, EventFlags: units2flags(units) | os.EventNoteAbsoluteTime, UserData: userData}
	return os.RegisterEventsWithQueue(ctx, q.KernelQueue, events)
}

func (q *Queue) Close(ctx *context.Context) bool {
	return os.CloseHandle(ctx, q.KernelQueue)
}

func (q *Queue) GetEvents(ctx *context.Context, events []os.Event) (int, bool) {
	if q.head < q.tail {
		n := copy(events, q.events[q.head:q.tail])
		q.head += n
		if q.head == q.tail {
			q.head = 0
			q.tail = 0
		}
		return n, true
	}
	return os.ReturnPendingEventsFromQueue(ctx, q.KernelQueue, events, nil)
}

func (q *Queue) requestNewEvents(ctx *context.Context, t *os.SecondsWithNanoseconds) bool {
	n, ok := os.ReturnPendingEventsFromQueue(ctx, q.KernelQueue, q.events[:], t)
	if !ok {
		return false
	}

	q.tail += n
	return true
}

func (q *Queue) HasEvents(ctx *context.Context) bool {
	if q.head < q.tail {
		return true
	}

	var t os.SecondsWithNanoseconds
	if !q.requestNewEvents(ctx, &t) {
		return false
	}

	return q.tail > 0
}

func (q *Queue) SyncFPS(ctx *context.Context, fps int) {
	now := cpu.ReadPerformanceCounter()
	durationBetweenPauses := now - q.LastSync
	targetRate := int64(float64(time.Second/time.Millisecond) / float64(fps) * float64(time.Millisecond))

	duration := targetRate - durationBetweenPauses.ToNanosecondsTruncated()
	if duration > 0 {
		if q.head >= q.tail {
			t := os.SecondsWithNanoseconds{Seconds: int(duration / time.Second), Nanoseconds: int(duration % time.Second)}
			q.requestNewEvents(ctx, &t)
		}
		now = cpu.ReadPerformanceCounter()
	}
	q.LastSync = now
}
