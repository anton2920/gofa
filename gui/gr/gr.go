package gr

import (
	"github.com/anton2920/gofa/gui/color"
	"github.com/anton2920/gofa/util"
)

var blacktext [256]color.Color

func init() {
	for i := 0; i < len(blacktext); i++ {
		blacktext[i] = color.Color(0x00010101*(255-i) + 0xff000000)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

/* TODO(anton2920): this is slow!!! */
func DrawCircle(pixmap Pixmap, bounds Rect, x0, y0, radius int, clr color.Color) {
	x1 := x0 + radius + radius
	y1 := y0 + radius + radius

	if x0 >= bounds.X1 {
		return
	}
	if x1 < bounds.X0 {
		return
	}
	if y0 >= bounds.Y1 {
		return
	}
	if y1 < bounds.Y0 {
		return
	}
	if clr.Invisible() {
		return
	}

	cx := x0 + radius
	cy := y0 + radius

	x0 = util.Max(x0, bounds.X0)
	x1 = util.Min(x1, bounds.X1-1)

	y0 = util.Max(y0, bounds.Y0)
	y1 = util.Min(y1, bounds.Y1-1)

	if clr.Opaque() {
		for y := y0; y < y1; y++ {
			for x := x0; x < x1; x++ {
				if (x-cx)*(x-cx)+(y-cy)*(y-cy) <= radius*radius {
					offset := y*pixmap.Stride + x
					pixmap.Pixels[offset] = clr
				}
			}
		}
	} else {
		for y := y0; y < y1; y++ {
			for x := x0; x < x1; x++ {
				if (x-cx)*(x-cx)+(y-cy)*(y-cy) <= radius*radius {
					offset := y*pixmap.Stride + x
					pixmap.Pixels[offset] = color.Blend(pixmap.Pixels[offset], clr)
				}
			}
		}

	}
}

func DrawLine(pixmap Pixmap, bounds Rect, x0, y0, x1, y1 int, clr color.Color) {
	if x0 == x1 {
		DrawVLine(pixmap, bounds, x0, y0, y1, clr)
	} else if y0 == y1 {
		DrawHLine(pixmap, bounds, y0, x0, x1, clr)
	} else {
		/* Not the fastest way to draw a line, but who wants to write 8 cases? */
		dx := x1 - x0
		dy := y1 - y0
		x := x0
		y := y0

		var tmp Rect
		tmp.X0 = util.Min(x0, x1)
		tmp.Y0 = util.Min(y0, y1)
		tmp.X1 = util.Max(x0, x1) + 1
		tmp.Y1 = util.Max(y0, y1) + 1

		if (tmp.X0 >= bounds.X1) || (tmp.X1 <= bounds.X0) {
			return
		}
		if (tmp.Y0 >= bounds.Y1) || (tmp.Y1 <= bounds.Y0) {
			return
		}
		if clr.Invisible() {
			return
		}

		opaque := clr.Opaque()
		maxLen := util.Max(abs(dx), abs(dy))
		invLen := float32(0xFFFF) / float32(maxLen)

		dx = int(float32(dx) * invLen)
		dy = int(float32(dy) * invLen)
		x = (x << 16) + 0x8000
		y = (y << 16) + 0x8000

		/* Does this line need clipping? */
		if bounds.Contains(tmp) {
			/* This is a very slow and dumb way to clip! */
			for i := 0; i <= maxLen; {
				if ((x >> 16) >= bounds.X0) && ((x >> 16) < bounds.X1) && ((y >> 16) >= bounds.Y0) && ((y >> 16) < bounds.Y1) {
					offset := (y>>16)*pixmap.Stride + (x >> 16)
					if opaque {
						pixmap.Pixels[offset] = clr
					} else {
						pixmap.Pixels[offset] = color.Blend(pixmap.Pixels[offset], clr)
					}
				}

				x += dx
				y += dy
				i++
			}
		} else if opaque {
			for i := 0; i <= maxLen; {
				offset := (y>>16)*pixmap.Stride + (x >> 16)
				pixmap.Pixels[offset] = clr
				x += dx
				y += dy
				i++
			}
		} else {
			for i := 0; i <= maxLen; {
				offset := (y>>16)*pixmap.Stride + (x >> 16)
				pixmap.Pixels[offset] = color.Blend(pixmap.Pixels[offset], clr)
				x += dx
				y += dy
				i++
			}
		}
	}
}

func DrawHLine(pixmap Pixmap, bounds Rect, y, x0, x1 int, clr color.Color) {
	if (y < bounds.Y0) || (y >= bounds.Y1) {
		return
	}
	if x1 < x0 {
		x0, x1 = x1, x0
	}
	if (x0 >= bounds.X1) || (x1 < bounds.X0) {
		return
	}
	if clr.Invisible() {
		return
	}

	x0 = util.Max(x0, bounds.X0)
	x1 = util.Min(x1, bounds.X1-1)

	n := x1 - x0 + 1

	out := pixmap.Pixels[y*pixmap.Stride+x0:]
	if clr.Opaque() {
		for i := 0; i < n; i++ {
			out[i] = clr
		}
	} else {
		for i := 0; i < n; i++ {
			out[i] = color.Blend(out[i], clr)
		}
	}
}

func DrawVLine(pixmap Pixmap, bounds Rect, x, y0, y1 int, clr color.Color) {
	if (x < bounds.X0) || (x >= bounds.X1) {
		return
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	if (y0 >= bounds.Y1) || (y1 < bounds.Y0) {
		return
	}
	if clr.Invisible() {
		return
	}

	y0 = util.Max(y0, bounds.Y0)
	y1 = util.Min(y1, bounds.Y1-1)

	n := y1 - y0 + 1

	out := pixmap.Pixels[y0*pixmap.Stride+x:]
	if clr.Opaque() {
		for i := 0; i < n; i++ {
			out[i*pixmap.Stride] = clr
		}
	} else {
		for i := 0; i < n; i++ {
			out[i*pixmap.Stride] = color.Blend(out[i*pixmap.Stride], clr)
		}
	}
}

func advanceBy(vs []color.Color, by int) []color.Color {
	if len(vs) < by {
		return nil
	}
	return vs[by:]
}

func drawPixmap(dst Pixmap, bounds Rect, x, y int, src Pixmap, pclr *color.Color) {
	srcBox := Rect{x, y, x + src.Width, y + src.Height}

	if !bounds.Contains(srcBox) {
		// runtime.Breakpoint()

		var x0, y0 int
		x1 := src.Width
		y1 := src.Height

		if bounds.X0 >= srcBox.X1 {
			return
		}
		if bounds.X1 <= srcBox.X0 {
			return
		}
		if bounds.Y0 >= srcBox.Y1 {
			return
		}
		if bounds.Y1 <= srcBox.Y0 {
			return
		}

		if x < bounds.X0 {
			x0 = bounds.X0 - x
		}
		if x+src.Width > bounds.X1 {
			x1 = bounds.X1 - x
		}
		if y < bounds.Y0 {
			y0 = bounds.Y0 - y
		}
		if y+src.Height > bounds.Y1 {
			y1 = bounds.Y1 - y
		}

		src = src.Sub(x0, y0, x1, y1)
		if src.Width <= 0 {
			panic("src.Width <= 0")
		}
		if src.Height <= 0 {
			panic("src.Height <= 0")
		}

		x += x0
		y += y0
	}

	/* Now the bitmap is clipped to be strictly onscreen. */
	out := dst.Pixels[y*dst.Stride+x:]
	in := src.Pixels
	if pclr != nil {
		clr := *pclr
		if (src.Alpha == AlphaFont) && (clr.Opaque()) {
			for j := 0; j < src.Height; j++ {
				if clr == color.Black {
					for i := 0; i < src.Width; i++ {
						if !in[i].Invisible() {
							if out[i] == color.White {
								out[i] = blacktext[in[i].A()]
							} else {
								out[i] = color.BlendMultiplyFont(out[i], in[i], clr)
							}
						}
					}
				} else {
					for i := 0; i < src.Width; i++ {
						if !in[i].Invisible() {
							out[i] = color.BlendMultiplyFont(out[i], in[i], clr)
						}
					}
				}

				out = advanceBy(out, dst.Stride)
				in = advanceBy(in, src.Stride)
			}
		} else {
			for j := 0; j < src.Height; j++ {
				for i := 0; i < src.Width; i++ {
					if !in[i].Invisible() {
						out[i] = color.BlendMultiply(out[i], in[i], clr)
					}
				}
				out = advanceBy(out, dst.Stride)
				in = advanceBy(in, src.Stride)
			}
		}
	} else if src.Alpha == AlphaOpaque {
		for j := 0; j < src.Height; j++ {
			copy(out, in[:src.Width])
			out = advanceBy(out, dst.Stride)
			in = advanceBy(in, src.Stride)
		}
	} else if src.Alpha == Alpha1bit {
		for j := 0; j < src.Height; j++ {
			for i := 0; i < src.Width; i++ {
				if !in[i].Invisible() {
					out[i] = in[i]
				}
			}
			out = advanceBy(out, dst.Stride)
			in = advanceBy(in, src.Stride)
		}
	} else {
		for j := 0; j < src.Height; j++ {
			for i := 0; i < src.Width; i++ {
				if !in[i].Invisible() {
					out[i] = color.Blend(out[i], in[i])
				}
			}
			out = advanceBy(out, dst.Stride)
			in = advanceBy(in, src.Stride)
		}
	}
}

func DrawPixmap(pixmap Pixmap, bounds Rect, x, y int, src Pixmap) {
	drawPixmap(pixmap, bounds, x, y, src, nil)
}

func DrawPixmapColored(pixmap Pixmap, bounds Rect, x, y int, src Pixmap, clr color.Color) {
	var pclr *color.Color
	if clr != 0xFFFFFFFF {
		pclr = &clr
	}
	drawPixmap(pixmap, bounds, x, y, src, pclr)
}

func DrawPoint(pixmap Pixmap, bounds Rect, x, y int, clr color.Color) {
	if (x >= bounds.X0) && (x < bounds.X1) && (y >= bounds.Y0) && (y < bounds.Y1) {
		offset := y*pixmap.Stride + x
		if clr.Opaque() {
			pixmap.Pixels[offset] = clr
		} else if !clr.Invisible() {
			pixmap.Pixels[offset] = color.Blend(pixmap.Pixels[offset], clr)
		}
	}
}

func DrawRectOutline(pixmap Pixmap, bounds Rect, x0, y0, x1, y1 int, clr color.Color) {
	if x1 < x0 {
		x0, x1 = x1, x0
	}
	if y1 < y0 {
		y0, y1 = y1, y0
	}
	DrawRectOutlineWH(pixmap, bounds, x0, y0, x1-x0+1, y1-y0+1, clr)
}

func DrawRectOutlineWH(pixmap Pixmap, bounds Rect, x0, y0, width, height int, clr color.Color) {
	if height == 1 {
		DrawHLine(pixmap, bounds, y0, x0, x0+width-1, clr)
	} else if width == 1 {
		DrawVLine(pixmap, bounds, x0, y0, y0+height-1, clr)
	} else if (height > 1) && (width > 1) {
		x1 := x0 + width - 1
		y1 := y0 + height - 1
		DrawHLine(pixmap, bounds, y0, x0, x1-1, clr)
		DrawVLine(pixmap, bounds, x1, y0, y1-1, clr)
		DrawHLine(pixmap, bounds, y1, x0+1, x1, clr)
		DrawVLine(pixmap, bounds, x0, y0+1, y1, clr)
	}
}

func DrawRectSolid(pixmap Pixmap, bounds Rect, x0, y0, x1, y1 int, clr color.Color) {
	if x1 < x0 {
		x0, x1 = x1, x0
	}
	if y1 < y0 {
		y0, y1 = y1, y0
	}
	DrawRectSolidWH(pixmap, bounds, x0, y0, x1-x0+1, y1-y0+1, clr)
}

func DrawRectSolidWH(pixmap Pixmap, bounds Rect, x0, y0, width, height int, clr color.Color) {
	if width > 0 {
		x1 := x0 + width - 1
		for j := 0; j < height; j++ {
			DrawHLine(pixmap, bounds, y0+j, x0, x1, clr)
		}
	}
}

func DrawText(pixmap Pixmap, bounds Rect, text string, font Font, x, y int, clr color.Color) {
	for i := 0; i < len(text); {
		c := text[i]
		i++

		if (c < font.startChar) || (c >= font.startChar+byte(len(font.fontChars))) {
			c = font.startChar
		}
		DrawPixmapColored(pixmap, bounds, x, y, font.fontChars[c-font.startChar], clr)
		x += font.fontChars[c-font.startChar].Width
		if (text[i-1] == 'f') && (text[i] == 't') {
			x--
		}
	}
}
