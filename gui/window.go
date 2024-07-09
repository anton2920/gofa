package gui

import (
	"github.com/anton2920/gofa/intel"
	"github.com/anton2920/gofa/time"
)

type WindowFlags uint

const (
	WindowNone WindowFlags = iota << 1
	WindowHidden
	WindowResizable
	WindowMinimized
	WindowMaximized
)

type Window struct {
	platformWindow

	Title  string
	Width  int
	Height int
	Flags  WindowFlags

	LastSync intel.Cycles
}

func NewWindow(title string, width, height int, flags WindowFlags) (*Window, error) {
	w := Window{Title: title, Width: width, Height: height, Flags: flags}

	if err := platformNewWindow(&w); err != nil {
		return nil, err
	}

	return &w, nil
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
