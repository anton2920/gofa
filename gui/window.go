package gui

import (
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/cpu"
	"github.com/anton2920/gofa/gui/color"
	"github.com/anton2920/gofa/time"
	"github.com/anton2920/gofa/time/time_"
	"github.com/anton2920/gofa/trace/trace_"
)

type WindowFlags uint

const (
	WindowNone WindowFlags = 1 << iota
	WindowHidden
	WindowResizable
	windowTransient
)

type Window struct {
	platformWindow
	Parent *Window

	Width  int
	Height int
	Flags  WindowFlags

	LastSyncCycles cpu.Cycles

	FrameTime float32
	FPS       float32

	title         [256]byte
	CursorVisible bool
}

func (w *Window) Init(title string, width int, height int, flags WindowFlags) error {
	w.Width = width
	w.Height = height
	w.Flags = flags
	w.CursorVisible = true

	if err := platformNewWindow(w, 0, 0); err != nil {
		return err
	}
	w.SetTitle(title)

	return nil
}

func (w *Window) Title() string {
	for i := 0; i < len(w.title); i++ {
		if w.title[i] == 0 {
			return bytes.AsString(w.title[:i])
		}
	}
	return ""
}

func (w *Window) SetTitle(title string) {
	w.title[copy(w.title[:len(w.title)-1], title)] = 0
	platformWindowSetTitle(w, &w.title[0])
}

func (w *Window) HasEvents() bool {
	t := trace_.Begin("")

	has := platformWindowHasEvents(w)

	trace_.End(t)
	return has
}

func (w *Window) GetEvents(events []Event) (int, error) {
	t := trace_.Begin("")

	n, err := platformWindowGetEvents(w, events)

	trace_.End(t)
	return n, err
}

func (w *Window) Invalidate() {
	t := trace_.Begin("")

	platformWindowInvalidate(w)

	trace_.End(t)
}

func (w *Window) DisplayPixels(pixels []color.Color, width, height int) {
	t := trace_.Begin("")

	platformWindowDisplayPixels(w, pixels, width, height)

	trace_.End(t)
}

func (w *Window) ShowCursor() {
	t := trace_.Begin("")

	if !w.CursorVisible {
		platformWindowEnableCursor(w)
		w.CursorVisible = true
	}

	trace_.End(t)
}

func (w *Window) HideCursor() {
	t := trace_.Begin("")

	if w.CursorVisible {
		platformWindowDisableCursor(w)
		w.CursorVisible = false
	}

	trace_.End(t)
}

func (w *Window) SyncFPS(fps int) {
	t := trace_.Begin("")

	now := cpu.ReadPerformanceCounter()
	durationBetweenPauses := (now - w.LastSyncCycles).ToNanosecondsTruncated()

	var targetRate int64
	if fps > 0 {
		targetRate = int64(float64(time.Second/time.Millisecond) / float64(fps) * float64(time.Millisecond))
	}
	dt := durationBetweenPauses

	duration := targetRate - durationBetweenPauses
	if duration > 0 {
		dt += duration
		time_.Sleep(duration)
		now = cpu.ReadPerformanceCounter()
	}

	if w.LastSyncCycles == 0 {
		dt = 0
	}
	w.LastSyncCycles = now

	w.FrameTime = float32(dt) / float32(time.Second)
	w.FPS = 1 / w.FrameTime

	// fmt.Printf("[gui]: Between: %d, Pause: %d, FrameTime: %g, FPS: %g\n", durationBetweenPauses, duration, w.FrameTime, w.FPS)

	trace_.End(t)
}

func (w *Window) Close() {
	platformWindowClose(w)
}
