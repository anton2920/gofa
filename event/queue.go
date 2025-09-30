package event

import (
	"unsafe"

	"github.com/anton2920/gofa/cpu"
	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/time"
)

type Queue struct {
	platformEventQueue

	LastSync cpu.Cycles
}

type Request int

const (
	RequestNone = Request(1 << iota)
	RequestRead
	RequestWrite
	RequestCount
)

type Trigger int

const (
	TriggerNone = Trigger(iota)
	TriggerLevel
	TriggerEdge
	TriggerCount
)

func NewQueue() (*Queue, error) {
	q := new(Queue)
	if err := platformNewEventQueue(q); err != nil {
		return nil, err
	}
	return q, nil
}

func (q *Queue) AddSocket(sock int32, request Request, trigger Trigger, userData unsafe.Pointer) error {
	return platformQueueAddSocket(q, sock, request, trigger, userData)
}

func (q *Queue) AddSignals(sigs ...syscall.Signal) error {
	for i := 0; i < len(sigs); i++ {
		if err := platformQueueAddSignal(q, sigs[i]); err != nil {
			return err
		}
	}
	return nil
}

func (q *Queue) AddTimer(id uintptr, duration int64, userData unsafe.Pointer) error {
	return platformQueueAddTimer(q, id, duration, userData)
}

func (q *Queue) AppendEvent(event Event) {
	platformQueueAppendEvent(q, event)
}

func (q *Queue) Close() error {
	return platformQueueClose(q)
}

func (q *Queue) GetEvents(events []Event) (int, error) {
	return platformQueueGetEvents(q, events)
}

func (q *Queue) HasEvents() bool {
	return platformQueueHasEvents(q)
}

func (q *Queue) SyncFPS(fps int) {
	now := cpu.ReadPerformanceCounter()
	durationBetweenPauses := now - q.LastSync
	targetRate := int64(1000 / float64(fps) * float64(time.Millisecond))

	duration := targetRate - durationBetweenPauses.ToNanoseconds()
	if duration > 0 {
		platformQueuePause(q, duration)
		now = cpu.ReadPerformanceCounter()
	}
	q.LastSync = now
}
