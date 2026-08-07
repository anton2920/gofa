package os

func RegisterEventsWithQueue(q Handle, chlist []Event) error {
	_, err := RegisterAndReturnPendingEventsFromQueue(q, chlist, nil, nil)
	return err
}

func ReturnPendingEventsFromQueue(q Handle, evlist []Event, t *SecondsWithNanoseconds) (int, error) {
	return RegisterAndReturnPendingEventsFromQueue(q, nil, evlist, t)
}
