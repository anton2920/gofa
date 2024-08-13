package gr

import (
	"github.com/anton2920/gofa/gui/color"
	"github.com/anton2920/gofa/util"
)

type Font struct {
	fontChars []Pixmap
	startChar byte
}

func skipBit(font *[]uint32, buffer, bitsLeft *int) {
	*buffer >>= 1
	*bitsLeft--

	if *bitsLeft == 0 {
		*buffer = int((*font)[0])
		*font = (*font)[1:]
		*bitsLeft = 32
	}
}

func DecompressFont(font []uint32) Font {
	var result Font

	start := byte((font[0] >> 0) & 0xFF)
	count := int((font[0] >> 8) & 0xFF)
	height := int((font[0] >> 16) & 0xFF)

	chars := make([]Pixmap, count)
	result.fontChars = chars
	result.startChar = start

	font = font[1:]
	for i := uint(0); i < uint(len(chars)); i++ {
		width := int((font[i>>2] >> ((i & 3) << 3)) & 0xFF)
		chars[i] = NewPixmap(width, height, AlphaFont)
	}
	font = font[(len(chars)+3)>>2:]

	buffer := int(font[0])
	font = font[1:]

	bitsLeft := 32
	for k := 0; k < len(chars); k++ {
		c := chars[k].Pixels

		for j := 0; j < chars[k].Height; j++ {
			z := buffer & 1
			skipBit(&font, &buffer, &bitsLeft)
			if z == 0 {
				for i := 0; i < chars[k].Width; i++ {
					c[j*chars[k].Width+i] = color.RGBA(255, 255, 255, 0)
				}
			} else {
				for i := 0; i < chars[k].Width; i++ {
					z = buffer & 1
					skipBit(&font, &buffer, &bitsLeft)
					if z == 0 {
						c[j*chars[k].Width+i] = color.RGBA(255, 255, 255, 0)
					} else {
						n := 0
						n += n + (buffer & 1)
						skipBit(&font, &buffer, &bitsLeft)
						n += n + (buffer & 1)
						skipBit(&font, &buffer, &bitsLeft)
						n += n + (buffer & 1)
						skipBit(&font, &buffer, &bitsLeft)
						n += 1
						c[j*chars[k].Width+i] = color.RGBA(255, 255, 255, byte((255*n)>>3))
					}
				}
			}
		}
	}

	return result
}

func (f *Font) CharHeight(c byte) int {
	if (c < f.startChar) || (c >= f.startChar+byte(len(f.fontChars))) {
		c = f.startChar
	}
	return f.fontChars[c-f.startChar].Height
}

func (f *Font) CharWidth(c byte) int {
	if (c < f.startChar) || (c >= f.startChar+byte(len(f.fontChars))) {
		c = f.startChar
	}
	return f.fontChars[c-f.startChar].Width
}

func (f *Font) TextHeight(text string) int {
	var height int

	for i := 0; i < len(text); i++ {
		height = util.Max(height, f.CharHeight(text[i]))
	}

	return height
}

func (f *Font) TextWidth(text string) int {
	var width int
	var i int

	for i < len(text) {
		width += f.CharWidth(text[i])
		i++
		if (text[i-1] == 'f') && (text[i] == 't') {
			width--
		}
	}

	return width
}
