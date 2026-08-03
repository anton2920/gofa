//go:build freebsd
// +build freebsd

package os

import (
	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/os/posix/freebsd"
)

type Event struct{}

func CreateNewEventQueue() (Handle, error) {
	q, err := freebsd.Kqueue()
	return Handle(q), err
}

func RegisterEventsWithEventQueue(q Handle, chlist []Event) error {
	return errors.ErrNotImplemented
}

func ReturnPendingEventsFromEventQueue(q Handle, evlist []Event, tv int64) (int, error) {
	return -1, errors.ErrNotImplemented
}
