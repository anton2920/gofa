package gr

import (
	"image"
	"reflect"
	"unsafe"

	"github.com/anton2920/gofa/gui/color"
)

type AlphaType int

const (
	AlphaOpaque AlphaType = iota
	Alpha1bit
	Alpha8bit
	AlphaFont
)

type Pixmap struct {
	Pixels []color.Color
	Width  int
	Height int
	Stride int
	Alpha  AlphaType
}

func NewPixmap(width, height int, alpha AlphaType) Pixmap {
	var pixmap Pixmap

	pixmap.Pixels = make([]color.Color, width*height)
	pixmap.Width = width
	pixmap.Height = height
	pixmap.Stride = width
	pixmap.Alpha = alpha

	/* Force alpha to be opaque. */
	for i := 0; i < len(pixmap.Pixels); i++ {
		pixmap.Pixels[i] = color.Black
	}
	return pixmap
}

func NewPixmapFromImage(img *image.RGBA, alpha AlphaType) Pixmap {
	var pixmap Pixmap

	pixmap.Pixels = *(*[]color.Color)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&img.Pix[0])), Len: len(img.Pix) >> 2, Cap: len(img.Pix) >> 2}))
	pixmap.Width = img.Bounds().Max.X - img.Bounds().Min.X
	pixmap.Height = img.Bounds().Max.Y - img.Bounds().Min.Y
	pixmap.Stride = pixmap.Width
	pixmap.Alpha = alpha

	for i := 0; i < len(pixmap.Pixels); i++ {
		clr := pixmap.Pixels[i]
		r := clr.B()
		g := clr.G()
		b := clr.R()
		a := clr.A()
		pixmap.Pixels[i] = color.RGBA(r, g, b, a)
	}

	return pixmap
}

func (p Pixmap) Sub(x0, y0, x1, y1 int) Pixmap {
	if x0 > x1 {
		panic("x0 > x1")
	}
	if y0 > y1 {
		panic("y0 > y1")
	}
	if x0 < 0 {
		x0 = 0
	}
	if y0 < 0 {
		y0 = 0
	}
	if x1 > p.Width {
		x1 = p.Width
	}
	if y1 > p.Height {
		y1 = p.Height
	}

	var out Pixmap
	out.Pixels = p.Pixels[x0+y0*p.Stride:]
	out.Width = x1 - x0
	out.Height = y1 - y0
	out.Stride = p.Stride
	out.Alpha = p.Alpha
	return out
}

func (p Pixmap) SubWH(x, y, width, height int) Pixmap {
	return p.Sub(x, y, x+width, y+height)
}
