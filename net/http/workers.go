package http

import (
	"fmt"

	"github.com/anton2920/gofa/event"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/trace"
)

type Workers struct {
	Queues  []*event.Queue
	Current int
}

func Worker(q *event.Queue, router Router) {
	rs := make([]Request, Pipeline)
	ws := make([]Response, Pipeline)

	events := make([]event.Event, 64)
	for {
		n, err := q.GetEvents(events)
		if err != nil {
			log.Errorf("Failed to read events: %v", err)
		}

		for i := 0; i < n; i++ {
			e := &events[i]
			if errno := e.Error(); errno != 0 {
				log.Errorf("Event for %d returned code %d (%s)", e.Identifier, errno, errno)
				continue
			}

			c, ok := ConnFromPointer(e.UserData)
			if !ok {
				continue
			}
			if e.EndOfFile() {
				c.Close()
				continue
			}

			switch e.Type {
			case event.TypeRead:
				n, err = c.ReadRequests(rs)
				if err != nil {
					if err.(syscall.Error).Errno != syscall.EAGAIN {
						continue
					}
					log.Errorf("Failed to read HTTP requests: %v", err)
					c.Close()
					continue
				}
				if n > 0 {
					RequestsHandler(ws[:n], rs[:n], router)
					FillResponses(c, ws[:n])
				}
				fallthrough
			case event.TypeWrite:
				if _, err := c.WriteResponses(nil); err != nil {
					log.Errorf("Failed to write HTTP responses: %v", err)
					c.Close()
					continue
				}
			}
		}
	}
}

func NewWorkers(router Router, n int) (*Workers, error) {
	var ws Workers
	var err error

	ws.Queues = make([]*event.Queue, ints.Max(1, n))
	for i := 0; i < len(ws.Queues); i++ {
		ws.Queues[i], err = event.NewQueue()
		if err != nil {
			return nil, fmt.Errorf("failed to create new event queue #%d: %v", i, err)
		}
		go Worker(ws.Queues[i], router)
	}

	return &ws, nil
}

func (ws *Workers) Add(c *Conn) error {
	t := trace.Begin("")

	/* TODO(anton2920): remove syscall! */
	flags, err := syscall.Fcntl(int32(c.Socket), syscall.F_GETFL, 0)
	if err != nil {
		trace.End(t)
		return fmt.Errorf("failed to get connection flags: %v", err)
	}
	if flags&syscall.O_NONBLOCK == 0 {
		flags |= syscall.O_NONBLOCK
		if _, err := syscall.Fcntl(int32(c.Socket), syscall.F_SETFL, flags); err != nil {
			trace.End(t)
			return fmt.Errorf("failed to set connection to non-blocking: %v", err)
		}
	}

	err = ws.Queues[ws.Current].AddSocket(int32(c.Socket), event.RequestRead|event.RequestWrite, event.TriggerEdge, c.Pointer())
	ws.Current = (ws.Current + 1) % len(ws.Queues)

	trace.End(t)
	return err
}
