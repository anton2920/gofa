/* Everything here is backported from 'strconv'. */
package fmt

func quotedStringLength(s string) int {
	var runeTmp [utfMax]byte
	quote := '"'
	size := 2

	for width := 0; len(s) > 0; s = s[width:] {
		r := rune(s[0])
		width = 1
		if r >= runeSelf {
			r, width = decodeRuneInString(s)
		}
		if (width == 1) && (r == runeError) {
			size += 2 + 2
			continue
		}
		if (r == rune(quote)) || (r == '\\') {
			size += 2
			continue
		}
		if isPrint(r) {
			size += encodeRune(runeTmp[:], r)
			continue
		}
		switch r {
		case '\a', '\b', '\f', '\n', '\r', '\t', '\v':
			size++
		default:
			switch {
			case r < ' ':
				size += 2 + 2
			case r > maxRune:
				r = 0xFFFD
				fallthrough
			case r < 0x10000:
				size += 2 + 4
			default:
				size += 2 + 8
			}
		}
	}

	return size
}

var lowerhex = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}

func (f *Formatter) quoteString(s string) {
	var runeTmp [utfMax]byte
	quote := '"'

	f.Pos += copy(f.Buffer[f.Pos:], "\"")
	for width := 0; len(s) > 0; s = s[width:] {
		r := rune(s[0])
		width = 1
		if r >= runeSelf {
			r, width = decodeRuneInString(s)
		}
		if width == 1 && r == runeError {
			f.Pos += copy(f.Buffer[f.Pos:], `\x`)
			f.Pos += copy(f.Buffer[f.Pos:], lowerhex[s[0]>>4])
			f.Pos += copy(f.Buffer[f.Pos:], lowerhex[s[0]&0xF])
			continue
		}
		if r == rune(quote) || r == '\\' {
			f.Pos += copy(f.Buffer[f.Pos:], "\\")
			f.Pos += copy(f.Buffer[f.Pos:], s[:1])
			continue
		}
		if isPrint(r) {
			n := encodeRune(runeTmp[:], r)
			f.Pos += copy(f.Buffer[f.Pos:], runeTmp[:n])
			continue
		}
		switch r {
		case '\a':
			f.Pos += copy(f.Buffer[f.Pos:], `\a`)
		case '\b':
			f.Pos += copy(f.Buffer[f.Pos:], `\b`)
		case '\f':
			f.Pos += copy(f.Buffer[f.Pos:], `\f`)
		case '\n':
			f.Pos += copy(f.Buffer[f.Pos:], `\n`)
		case '\r':
			f.Pos += copy(f.Buffer[f.Pos:], `\r`)
		case '\t':
			f.Pos += copy(f.Buffer[f.Pos:], `\t`)
		case '\v':
			f.Pos += copy(f.Buffer[f.Pos:], `\v`)
		default:
			switch {
			case r < ' ':
				f.Pos += copy(f.Buffer[f.Pos:], `\x`)
				f.Pos += copy(f.Buffer[f.Pos:], lowerhex[s[0]>>4])
				f.Pos += copy(f.Buffer[f.Pos:], lowerhex[s[0]&0xF])
			case r > maxRune:
				r = 0xFFFD
				fallthrough
			case r < 0x10000:
				f.Pos += copy(f.Buffer[f.Pos:], `\u`)
				for s := 12; s >= 0; s -= 4 {
					f.Pos += copy(f.Buffer[f.Pos:], lowerhex[r>>uint(s)&0xF])
				}
			default:
				f.Pos += copy(f.Buffer[f.Pos:], `\U`)
				for s := 28; s >= 0; s -= 4 {
					f.Pos += copy(f.Buffer[f.Pos:], lowerhex[r>>uint(s)&0xF])
				}
			}
		}
	}
	f.Pos += copy(f.Buffer[f.Pos:], "\"")
}
