package gui

import (
	"github.com/anton2920/gofa/intel"
	"github.com/anton2920/gofa/time"
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
	return platformWindowHasEvents(w)
}

func (w *Window) GetEvents(events []Event) (int, error) {
	return platformWindowGetEvents(w, events)
}

func (w *Window) Invalidate() {
	platformWindowInvalidate(w)
}

func (w *Window) DisplayPixels(pixels []uint32, width, height int) {
	platformWindowDisplayPixels(w, pixels, width, height)
}

func (w *Window) ShowCursor() {
	if !w.CursorVisible {
		platformWindowEnableCursor(w)
		w.CursorVisible = true
	}
}

func (w *Window) HideCursor() {
	if w.CursorVisible {
		platformWindowDisableCursor(w)
		w.CursorVisible = false
	}
}

func (w *Window) SyncFPS(fps int) {
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
}

func (w *Window) Close() {
	platformWindowClose(w)
}
