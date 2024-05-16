package event

import (
	"runtime"
	"unsafe"
)

type Queue struct {
	platformEventQueue

	Pinner    runtime.Pinner
	LastPause int64
}

type Request int

const (
	RequestRead Request = (1 << iota)
	RequestWrite
)

type Trigger int

const (
	TriggerLevel Trigger = iota
	TriggerEdge
)

type DurationMeasurement int

const (
	Seconds DurationMeasurement = iota
	Milliseconds
	Microseconds
	Nanoseconds
	Absolute
)

type Type int32

const (
	None Type = iota
	Read
	Write
	Signal
	Timer
)

type Event struct {
	Type       Type
	Identifier int32
	Available  int
	UserData   unsafe.Pointer

	/* TODO(anton2920): I don't like this!!! */
	EndOfFile bool
}

var EmptyEvent Event

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

func (q *Queue) AddSignal(sig int32) error {
	return platformQueueAddSignal(q, sig)
}

func (q *Queue) AddTimer(identifier int32, duration int, measurement DurationMeasurement, userData unsafe.Pointer) error {
	return platformQueueAddTimer(q, identifier, duration, measurement, userData)
}

func (q *Queue) AppendEvent(event Event) {
	platformQueueAppendEvent(q, event)
}

func (q *Queue) Close() error {
	q.Pinner.Unpin()
	return platformQueueClose(q)
}

func (q *Queue) GetEvent() (Event, error) {
	return platformQueueGetEvent(q)
}

func (q *Queue) HasEvents() bool {
	return platformQueueHasEvents(q)
}

func (q *Queue) Pause(FPS int) {
	now := platformQueueGetTime()
	durationBetweenPauses := now - q.LastPause
	targetRate := int64(1000.0/float32(FPS)) * 1_000_000

	duration := targetRate - durationBetweenPauses
	if duration > 0 {
		platformQueuePause(q, duration)
		now = platformQueueGetTime()
	}
	q.LastPause = now
}
