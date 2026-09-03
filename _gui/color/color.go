package color

type Color uint32

var (
	Black = RGB(0, 0, 0)
	Red   = RGB(255, 0, 0)
	Green = RGB(0, 255, 0)
	Blue  = RGB(0, 0, 255)
	White = RGB(255, 255, 255)
)

func (c Color) A() byte {
	return byte((c >> 24) & 0xFF)
}

func (c Color) B() byte {
	return byte(c & 0xFF)
}

func (c Color) G() byte {
	return byte((c >> 8) & 0xFF)
}

func (c Color) R() byte {
	return byte((c >> 16) & 0xFF)
}

func (c Color) Invisible() bool {
	return c <= 0x00FFFFFF
}

/* Opaque returns true if color is 100% non-transparent. */
func (c Color) Opaque() bool {
	return c >= 0xFF000000
}

func Blend(dst, src Color) Color {
	/* Accelerated blend computes r and b in parallel. */
	a := Color(src.A())
	rbSrc := src & 0xFF00FF
	rbDst := dst & 0xFF00FF
	rb := rbDst + ((rbSrc - rbDst) * a >> 8)
	gDst := dst & 0x00FF00
	g := gDst + (((src & 0x00FF00) - (dst & 0x00FF00)) * a >> 8)
	/* NOTE(anton2920): we do not compute a real dst alpha. */
	return (rb & 0xFF00FF) + (g & 0x00FF00) + 0xFF000000
}

func BlendMultiply(dst, src1, src2 Color) Color {
	sr := Color(src1.R()) * Color(src2.R()) >> 8
	sg := Color(src1.G()) * Color(src2.G()) >> 8
	sb := Color(src1.B()) * Color(src2.B()) >> 8
	sa := Color(src1.A()) * Color(src2.A()) >> 8

	r := sr + (sr >> 7) /* 0..255. */
	g := sg + (sg >> 7)
	b := sb + (sb >> 7)
	a := sa + ((sa >> 6) & 2) /* 0..256. */

	dr := r - Color(dst.R())
	dg := g - Color(dst.G())
	db := b - Color(dst.B())

	or := dst.R() + byte((dr*a)>>8)
	og := dst.G() + byte((dg*a)>>8)
	ob := dst.B() + byte((db*a)>>8)

	return RGB(or, og, ob)
}

func BlendMultiplyFont(dst, font, src Color) Color {
	/* Accelerated blend computes r and b in parallel. */
	a := Color(font.A())
	rbSrc := src & 0xFF00FF
	rbDst := dst & 0xFF00FF
	rb := rbDst + ((rbSrc - rbDst) * a >> 8)
	gDst := dst & 0x00FF00
	g := gDst + (((src & 0x00FF00) - (dst & 0x00FF00)) * a >> 8)
	/* NOTE(anton2920): we do not compute a real dst alpha. */
	return (rb & 0xFF00FF) + (g & 0x00FF00) + 0xFF000000
}

func Average(c, d Color) Color {
	c = (c >> 1) & 0x7f7f7f7f
	d = (d >> 1) & 0x7f7f7f7f
	return c + d + (c & d & 0x01010101)
}

func Dark(c Color) Color {
	return Average(c, Black)
}

func Grey(c Color) Color {
	return RGB(byte(c), byte(c), byte(c))
}

func Lite(c Color) Color {
	return Average(c, White)
}

func RGB(r, g, b byte) Color {
	return RGBA(r, g, b, 255)
}

func RGBA(r, g, b, a byte) Color {
	return Color(a)<<24 | Color(r)<<16 | Color(g)<<8 | Color(b)
}
