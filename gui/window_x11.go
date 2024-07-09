//go:build unix

package gui

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib
#cgo LDFLAGS: -lX11 -lm -lxcb -lXau -lXdmcp

#include <X11/Xlib.h>
#include <X11/Xutil.h>
*/
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/anton2920/gofa/errors"
)

type platformWindow struct {
	wmDeleteWindow C.Atom
	display        *C.Display
	window         C.Window
	root           C.Window
	visual         *C.Visual
	screen         C.int
	gc             C.GC

	pendingEvents int
}

func platformNewWindow(w *Window) error {
	w.display = C.XOpenDisplay(nil)
	if w.display == nil {
		return errors.New("failed to open display")
	}

	w.screen = C.XDefaultScreen(w.display)
	w.visual = C.XDefaultVisual(w.display, w.screen)
	if w.visual.class != C.TrueColor {
		return errors.New("cannot handle non-true color visual")
	}
	w.root = C.XDefaultRootWindow(w.display)
	w.gc = C.XDefaultGC(w.display, w.screen)

	w.window = C.XCreateSimpleWindow(w.display, w.root, 0, 0, C.uint(w.Width), C.uint(w.Height), 1, 0, 0)
	C.XSelectInput(w.display, w.window, C.ExposureMask|C.KeyPressMask|C.KeyReleaseMask|C.ButtonPressMask|C.ButtonReleaseMask|C.PointerMotionMask|C.StructureNotifyMask)

	platformWindowSetTitle(w, w.Title)

	if (w.Flags & WindowResizable) == 0 {
		hints := C.XAllocSizeHints()
		hints.flags = C.PMinSize | C.PMaxSize
		hints.min_width = C.int(w.Width)
		hints.min_height = C.int(w.Height)
		hints.max_width = C.int(w.Width)
		hints.max_height = C.int(w.Height)
		C.XSetWMNormalHints(w.display, w.window, hints)
		C.XFree(unsafe.Pointer(hints))
	}

	if (w.Flags & WindowHidden) == 0 {
		C.XMapWindow(w.display, w.window)
	}

	w.wmDeleteWindow = C.XInternAtom(w.display, C.CString("WM_DELETE_WINDOW"), 1)
	C.XSetWMProtocols(w.display, w.window, &w.wmDeleteWindow, 1)

	return nil
}

func platformWindowSetTitle(w *Window, title string) {
	C.XStoreName(w.display, w.window, C.CString(title))
}

func platformWindowHasEvents(w *Window) bool {
	w.pendingEvents = int(C.XPending(w.display))
	return w.pendingEvents > 0
}

func platformWindowGetEvents(w *Window, events []Event) (int, error) {
	var platformEvent C.XEvent

	n := min(w.pendingEvents, len(events))
	var consumed int
	for i := 0; i < n; i++ {
		event := &events[consumed]

		C.XNextEvent(w.display, &platformEvent)
		/* NOTE(anton2920): convoluted way of saying 'platformEvent.type'. */
		switch *(*C.int)(unsafe.Pointer(&platformEvent)) {
		case C.ClientMessage:
			clientEvent := *(*C.XClientMessageEvent)(unsafe.Pointer(&platformEvent))
			data := *(*C.int)(unsafe.Pointer(&clientEvent.data[0]))

			if C.Atom(data) == w.wmDeleteWindow {
				event.Type = DestroyEvent
				consumed++
			}
		case C.Expose:
			event.Type = PaintEvent
			consumed++
		case C.ConfigureNotify:
			configureEvent := *(*C.XConfigureEvent)(unsafe.Pointer(&platformEvent))
			eventWidth := int(configureEvent.width)
			eventHeight := int(configureEvent.height)

			if (eventWidth != w.Width) || (eventHeight != w.Height) {
				w.Width = eventWidth
				w.Height = eventHeight

				event.Type = ResizeEvent
				event.Width = eventWidth
				event.Height = eventHeight
				consumed++
			}
		case C.ButtonPress:
			buttonEvent := *(*C.XButtonEvent)(unsafe.Pointer(&platformEvent))
			eventX := int(buttonEvent.x)
			eventY := int(buttonEvent.y)

			event.Type = MousePressEvent
			event.Button = MouseButton(buttonEvent.button)
			event.X = eventX
			event.Y = eventY
			consumed++
		case C.ButtonRelease:
			buttonEvent := *(*C.XButtonEvent)(unsafe.Pointer(&platformEvent))
			eventX := int(buttonEvent.x)
			eventY := int(buttonEvent.y)

			event.Type = MouseReleaseEvent
			event.Button = MouseButton(buttonEvent.button)
			event.X = eventX
			event.Y = eventY
			consumed++
		case C.MotionNotify:
			motionEvent := *(*C.XMotionEvent)(unsafe.Pointer(&platformEvent))
			eventX := int(motionEvent.x)
			eventY := int(motionEvent.y)

			event.Type = MouseMoveEvent
			event.X = eventX
			event.Y = eventY
			consumed++
		}
	}

	return consumed, nil
}

func platformWindowInvalidate(w *Window) {
	C.XClearArea(w.display, w.window, 0, 0, 1, 1, C.int(1))
	C.XFlush(w.display)
}

func platformWindowDisplayPixels(w *Window, pixels []uint32, width int, height int) {
	var image C.XImage
	var pinner runtime.Pinner

	image.width = C.int(width)
	image.height = C.int(height)
	image.format = C.ZPixmap
	image.data = (*C.char)(unsafe.Pointer(unsafe.SliceData(pixels)))
	image.bitmap_unit = C.int(unsafe.Sizeof(pixels[0]) * 8)
	image.bitmap_pad = C.int(unsafe.Sizeof(pixels[0]) * 8)
	image.depth = 24
	image.bytes_per_line = C.int(width * int(unsafe.Sizeof(pixels[0])))
	image.bits_per_pixel = C.int(unsafe.Sizeof(pixels[0]) * 8)
	image.red_mask = w.visual.red_mask
	image.green_mask = w.visual.green_mask
	image.blue_mask = w.visual.blue_mask

	pinner.Pin(&pixels[0])
	C.XInitImage(&image)
	C.XPutImage(w.display, w.window, w.gc, &image, 0, 0, 0, 0, C.uint(width), C.uint(height))
	pinner.Unpin()
}

func platformWindowClose(w *Window) {
	C.XCloseDisplay(w.display)
}
