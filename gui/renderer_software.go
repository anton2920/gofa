package gui

import (
	"unsafe"

	"github.com/anton2920/gofa/gui/color"
	"github.com/anton2920/gofa/gui/gr"
	"github.com/anton2920/gofa/util"
)

/* SoftwareRenderer is a completely platform-independent renderer. All rendering is done in a Pixmap, which is then transferred to a screen via Window platform-specific function. This enables maximum portability at a cost of performance. */
type SoftwareRenderer struct {
	Window *Window
	Pixmap gr.Pixmap
	Active gr.Rect
}

func NewSoftwareRenderer(window *Window) *SoftwareRenderer {
	var r SoftwareRenderer

	r.Window = window
	r.Pixmap = gr.NewPixmap(window.Width, window.Height, gr.AlphaOpaque)
	r.Active = Rect{0, 0, r.Pixmap.Width, r.Pixmap.Height}

	return &r
}

func (r *SoftwareRenderer) Clear(clr color.Color) {
	util.Memset(r.Pixmap.Pixels, clr)
}

func (r *SoftwareRenderer) Present() {
	r.Window.DisplayPixels(unsafe.Slice((*uint32)(unsafe.Pointer(unsafe.SliceData(r.Pixmap.Pixels))), len(r.Pixmap.Pixels)), r.Pixmap.Width, r.Pixmap.Height)
}

func (r *SoftwareRenderer) Resize(width int, height int) {
	if (width > r.Pixmap.Width) || (height > r.Pixmap.Height) {
		r.Pixmap = gr.NewPixmap(width, height, gr.AlphaOpaque)
	} else {
		r.Pixmap.Pixels = r.Pixmap.Pixels[:width*height]
		r.Pixmap.Width = width
		r.Pixmap.Height = height
		r.Pixmap.Stride = width
	}
	r.Active = Rect{0, 0, r.Pixmap.Width, r.Pixmap.Height}
}

func (r *SoftwareRenderer) RenderPoint(x, y, size int, clr color.Color) {
	if size <= 1 {
		gr.DrawPoint(r.Pixmap, r.Active, x, y, clr)
	} else {
		gr.DrawRectSolid(r.Pixmap, r.Active, x-size, y-size, x+size, y+size, clr)
	}
}

func (r *SoftwareRenderer) RenderLine(x0, y0, x1, y1 int, clr color.Color) {
	gr.DrawLine(r.Pixmap, r.Active, x0, y0, x1, y1, clr)
}

func (r *SoftwareRenderer) RenderRect(x0, y0, x1, y1 int, clr color.Color) {
	gr.DrawRectOutline(r.Pixmap, r.Active, x0, y0, x1, y1, clr)
}

func (r *SoftwareRenderer) RenderRectWH(x, y, width, height int, clr color.Color) {
	gr.DrawRectOutlineWH(r.Pixmap, r.Active, x, y, width, height, clr)
}

func (r *SoftwareRenderer) RenderSolidRect(x0, y0, x1, y1 int, clr color.Color) {
	gr.DrawRectSolid(r.Pixmap, r.Active, x0, y0, x1, y1, clr)
}

func (r *SoftwareRenderer) RenderSolidRectWH(x, y, width, height int, clr color.Color) {
	gr.DrawRectSolidWH(r.Pixmap, r.Active, x, y, width, height, clr)
}

func (r *SoftwareRenderer) RenderCircle(x0, y0, radius int, clr color.Color) {
	gr.DrawCircle(r.Pixmap, r.Active, x0, y0, radius, clr)
}

func (r *SoftwareRenderer) RenderText(text string, font Font, x, y int, clr color.Color) {
	gr.DrawText(r.Pixmap, r.Active, text, font, x, y, clr)
}

func (r *SoftwareRenderer) RenderPixmap(pixmap gr.Pixmap, x, y int) {
	gr.DrawPixmap(r.Pixmap, r.Active, x, y, pixmap)
}
