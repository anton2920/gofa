package gui

import (
	"github.com/anton2920/gofa/gui/color"
	"github.com/anton2920/gofa/gui/gr"
)

type Renderer interface {
	/* Utils. */
	Clear(clr color.Color)
	Resize(width, height int)
	Present()

	/* 2D geometry. */
	RenderPoint(x, y, size int, clr color.Color)
	RenderLine(x0, y0, x1, y1 int, clr color.Color)
	RenderRect(x0, y0, x1, y1 int, clr color.Color)
	RenderRectWH(x, y, width, height int, clr color.Color)
	RenderSolidRect(x0, y0, x1, y1 int, clr color.Color)
	RenderSolidRectWH(x, y, width, height int, clr color.Color)
	RenderCircle(x0, y0, radius int, clr color.Color)

	/* Text. */
	RenderText(text string, font gr.Font, x, y int, clr color.Color)

	/* Graphics. */
	RenderPixmap(pixmap gr.Pixmap, x, y int)
}
