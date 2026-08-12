package os

import "github.com/anton2920/gofa/context"

func RegisterEventsWithQueue(ctx *context.Context, q Handle, chlist []Event) bool {
	_, ok := RegisterAndReturnPendingEventsFromQueue(ctx, q, chlist, nil, nil)
	return ok
}

func ReturnPendingEventsFromQueue(ctx *context.Context, q Handle, evlist []Event, t *SecondsWithNanoseconds) (int, bool) {
	return RegisterAndReturnPendingEventsFromQueue(ctx, q, nil, evlist, t)
}
