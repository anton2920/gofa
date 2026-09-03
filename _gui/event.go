package gui

type MouseButton int

const (
	Button1 MouseButton = iota + 1
	Button2
	Button3
	Button4
	Button5
)

type EventType int

const (
	DestroyEvent EventType = iota
	PaintEvent

	MouseMoveEvent
	MousePressEvent
	MouseReleaseEvent

	ResizeEvent
)

type Event struct {
	Type EventType

	/* For Mouse events. */
	X, Y   int
	Button MouseButton

	/* For Resize event. */
	Width, Height int
}
