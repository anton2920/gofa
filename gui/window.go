package gui

import (
	"github.com/anton2920/gofa/gui/color"
	"github.com/anton2920/gofa/intel"
	"github.com/anton2920/gofa/time"
	"github.com/anton2920/gofa/trace"
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

	Title  string
	Width  int
	Height int
	Flags  WindowFlags

	CursorVisible bool

	LastSync intel.Cycles
}

func NewWindow(title string, width, height int, flags WindowFlags) (*Window, error) {
	w := Window{Title: title, Width: width, Height: height, Flags: flags, CursorVisible: true}

	if err := platformNewWindow(&w, 0, 0); err != nil {
		return nil, err
	}

	return &w, nil
}

func (w *Window) NewTransientWindow(title string, x, y, width, height int) (*Window, error) {
	tw := Window{Parent: w, Title: title, Width: width, Height: height, Flags: windowTransient, CursorVisible: true}

	if err := platformNewWindow(&tw, x, y); err != nil {
		return nil, err
	}

	return &tw, nil
}

func (w *Window) SetTitle(title string) {
	w.Title = title
	platformWindowSetTitle(w, title)
}

func (w *Window) HasEvents() bool {
	t := trace.Begin("")

	has := platformWindowHasEvents(w)

	trace.End(t)
	return has
}

func (w *Window) GetEvents(events []Event) (int, error) {
	t := trace.Begin("")

	n, err := platformWindowGetEvents(w, events)

	trace.End(t)
	return n, err
}

func (w *Window) Invalidate() {
	t := trace.Begin("")

	platformWindowInvalidate(w)

	trace.End(t)
}

func (w *Window) DisplayPixels(pixels []color.Color, width, height int) {
	t := trace.Begin("")

	platformWindowDisplayPixels(w, pixels, width, height)

	trace.End(t)
}

func (w *Window) ShowCursor() {
	t := trace.Begin("")

	if !w.CursorVisible {
		platformWindowEnableCursor(w)
		w.CursorVisible = true
	}

	trace.End(t)
}

func (w *Window) HideCursor() {
	t := trace.Begin("")

	if w.CursorVisible {
		platformWindowDisableCursor(w)
		w.CursorVisible = false
	}

	trace.End(t)
}

func (w *Window) SyncFPS(fps int) {
	t := trace.Begin("")

	now := intel.RDTSC()
	durationBetweenPauses := now - w.LastSync
	targetRate := int64(time.MsecPerSec / float64(fps) * (time.NsecPerSec / time.MsecPerSec))

	duration := targetRate - durationBetweenPauses.ToNsec()
	if duration > 0 {
		platformSleep(duration)
		now = intel.RDTSC()
	}
	// println(int(time.MsecPerSec/float64(durationBetweenPauses.ToMsec())), "FPS")
	w.LastSync = now

	trace.End(t)
}

func (w *Window) Close() {
	platformWindowClose(w)
}
