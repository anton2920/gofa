package html

import (
	stdbytes "bytes"

	"github.com/anton2920/gofa/bytes"
)

func (h *HTML) Begin2() *HTML {
	return h.String(`<!DOCTYPE html><html>`)
}

func (h *HTML) End2() *HTML {
	return h.String(`</html>`)
}

func (h *HTML) HeadBegin2() *HTML {
	h.String(` <head>`)

	h.String(`<meta charset="UTF-8">`)
	h.String(`<meta name="viewport" content="width=device-width, initial-scale=1.0">`)

	if len(h.Theme.HeadLink.Href) > 0 {
		h.String(` <link href="`).String(h.Theme.HeadLink.Href).String(`" rel="`).String(h.Theme.HeadLink.Rel).String(`">`)
	}
	if len(h.Theme.HeadScript.Src) > 0 {
		h.String(` <script src="`).String(h.Theme.HeadScript.Src).String(`">`).String(`</script>`)
	}

	return h
}

func (h *HTML) HeadEnd2() *HTML {
	return h.String(`</head>`)
}

func (h *HTML) TitleBegin2() *HTML {
	return h.String(` <title>`)
}

func (h *HTML) Title2(title string) *HTML {
	return h.TitleBegin2().LString(title).TitleEnd2()
}

func (h *HTML) TitleEnd2() *HTML {
	return h.String(`</title>`)
}

func (h *HTML) BodyBegin2() *HTML {
	return h.String(` <body>`).Class_(h.Theme.Body.Class)
}

func (h *HTML) BodyEnd2() *HTML {
	return h.String(`</body>`)
}

func (h *HTML) DivBegin2() *HTML {
	return h.String(` <div>`)
}

func (h *HTML) DivEnd2() *HTML {
	return h.String(`</div>`)
}

func (h *HTML) UlBegin2() *HTML {
	return h.String(` <ul>`)
}

func (h *HTML) UlEnd2() *HTML {
	return h.String(`</ul>`)
}

func (h *HTML) OlBegin2() *HTML {
	return h.String(` <ol>`)
}

func (h *HTML) OlEnd2() *HTML {
	return h.String(`</ol>`)
}

func (h *HTML) LiBegin2() *HTML {
	return h.String(` <li>`)
}

func (h *HTML) LiEnd2() *HTML {
	return h.String(`</li>`)
}

func (h *HTML) PBegin2() *HTML {
	return h.String(` <p>`)
}

func (h *HTML) P2(p string) *HTML {
	return h.PBegin2().LString(p).PEnd2()
}

func (h *HTML) PEnd2() *HTML {
	return h.String(`</p>`)
}

func (h *HTML) BBegin2() *HTML {
	return h.String(` <b>`)
}

func (h *HTML) B2(b string) *HTML {
	return h.BBegin2().LString(b).BEnd2()
}

func (h *HTML) BEnd2() *HTML {
	return h.String(`</b>`)
}

func (h *HTML) IBegin2() *HTML {
	return h.String(` <i>`)
}

func (h *HTML) I2(i string) *HTML {
	return h.IBegin2().LString(i).IEnd2()
}

func (h *HTML) IEnd2() *HTML {
	return h.String(`</i>`)
}

func (h *HTML) SpanBegin2() *HTML {
	return h.String(` <span>`).Class_(h.Theme.Span.Class)
}

func (h *HTML) Span2(span string) *HTML {
	return h.SpanBegin2().LString(span).SpanEnd2()
}

func (h *HTML) SpanEnd2() *HTML {
	return h.String(`</span>`)
}

func (h *HTML) H1Begin2() *HTML {
	return h.String(` <h1>`).Class_(h.Theme.H1.Class)
}

func (h *HTML) H12(h1 string) *HTML {
	return h.H1Begin2().LString(h1).H1End2()
}

func (h *HTML) H1End2() *HTML {
	return h.String(`</h1>`)
}

func (h *HTML) H2Begin2() *HTML {
	return h.String(` <h2>`).Class_(h.Theme.H2.Class)
}

func (h *HTML) H22(h2 string) *HTML {
	return h.H2Begin2().LString(h2).H2End2()
}

func (h *HTML) H2End2() *HTML {
	return h.String(`</h2>`)
}

func (h *HTML) H3Begin2() *HTML {
	return h.String(` <h3>`).Class_(h.Theme.H3.Class)
}

func (h *HTML) H32(h3 string) *HTML {
	return h.H3Begin2().LString(h3).H3End2()
}

func (h *HTML) H3End2() *HTML {
	return h.String(`</h3>`)
}

func (h *HTML) H4Begin2() *HTML {
	return h.String(` <h4>`).Class_(h.Theme.H4.Class)
}

func (h *HTML) H42(h4 string) *HTML {
	return h.H4Begin2().LString(h4).H4End2()
}

func (h *HTML) H4End2() *HTML {
	return h.String(`</h4>`)
}

func (h *HTML) H5Begin2() *HTML {
	return h.String(` <h5>`).Class_(h.Theme.H5.Class)
}

func (h *HTML) H52(h5 string) *HTML {
	return h.H5Begin2().LString(h5).H5End2()
}

func (h *HTML) H5End2() *HTML {
	return h.String(`</h5>`)
}

func (h *HTML) H6Begin2() *HTML {
	return h.String(` <h6>`).Class_(h.Theme.H6.Class)
}

func (h *HTML) H62(h6 string) *HTML {
	return h.H6Begin2().LString(h6).H6End2()
}

func (h *HTML) H6End2() *HTML {
	return h.String(`</h6>`)
}

func (h *HTML) ABegin2(href string) *HTML {
	return h.String(` <a href="`).String(href).String(`">`).Class_(h.Theme.A.Class)
}

func (h *HTML) A2(href string, contents string) *HTML {
	return h.ABegin2(href).LString(contents).AEnd2()
}

func (h *HTML) AEnd2() *HTML {
	return h.String(`</a>`)
}

func (h *HTML) Link2(href string) *HTML {
	return h.String(`<link href="`).String(href).String(`">`)
}

func (h *HTML) Img2(alt string, src string) *HTML {
	return h.String(` <img alt="`).LString(alt).String(`" src="`).String(src).String(`">`)
}

func (h *HTML) Input2(typ string) *HTML {
	return h.String(` <input type="`).String(typ).String(`">`).Class_(h.Theme.Input.Class)
}

func (h *HTML) Button2(value string) *HTML {
	return h.WithoutTheme().Input2("submit").Value(value).Class_(h.Theme.Button.Class)
}

func (h *HTML) FormBegin2(method string) *HTML {
	return h.String(` <form method="`).String(method).String(`">`).Class_(h.Theme.Form.Class)
}

func (h *HTML) FormEnd2() *HTML {
	return h.String(`</form>`)
}

func (h *HTML) LabelBegin2() *HTML {
	return h.String(` <label>`).Class_(h.Theme.Label.Class)
}

func (h *HTML) Label2(label string) *HTML {
	return h.LabelBegin2().LString(label).LabelEnd2()
}

func (h *HTML) LabelEnd2() *HTML {
	return h.String(`</label>`)
}

func (h *HTML) Checkbox2() *HTML {
	return h.WithoutTheme().Input2("checkbox").Class_(h.Theme.Checkbox.Class)
}

func (h *HTML) SelectBegin2() *HTML {
	return h.String(` <select>`)
}

func (h *HTML) SelectEnd2() *HTML {
	return h.String(`</select>`)
}

func (h *HTML) OptionBegin2() *HTML {
	return h.String(` <option>`)
}

func (h *HTML) OptionEnd2() *HTML {
	return h.String(`</option>`)
}

func (h *HTML) SvgBegin2(width int, height int) *HTML {
	return h.String(` <svg xmlns="http://www.w3.org/2000/svg" width="`).Int(width).String(`" height="`).Int(height).String(`">`)
}

func (h *HTML) SvgEnd2() *HTML {
	return h.String(`</svg>`)
}

func (h *HTML) TextBegin2(x int, y int) *HTML {
	return h.String(` <text x="`).Int(x).String(`" y="`).Int(y).String(`">`)
}

func (h *HTML) TextEnd2() *HTML {
	return h.String(`</text>`)
}

func (h *HTML) Path2(d string) *HTML {
	return h.String(` <path d="`).String(d).String(`"/>`)
}

func (h *HTML) Circle2(cx int, cy int, r int) *HTML {
	return h.String(` <circle cx="`).Int(cx).String(`" cy="`).Int(cy).String(`" r="`).Int(r).String(`"/>`)
}

func (h *HTML) Rect2(x int, y int, width int, height int, rx int) *HTML {
	return h.String(` <rect x="`).Int(x).String(`" y="`).Int(y).String(`" width="`).Int(width).String(`" height="`).Int(height).String(`" rx="`).Int(rx).String(`"/>`)
}

func (h *HTML) Line2(x1 int, y1 int, x2 int, y2 int) *HTML {
	return h.String(` <line x1="`).Int(x1).String(`" y1="`).Int(y1).String(`" x2="`).Int(x2).String(`" y2="`).Int(y2).String(`"/>`)
}

var _class = []byte("class")

func (h *HTML) Class_(class string) *HTML {
	if len(class) > 0 {
		h.Backspace().String(` class="`).String(class).String(`">`)
	}
	return h
}

func (h *HTML) Class(class string) *HTML {
	if len(class) > 0 {
		closeTagBegin := -1

		/* TODO(anton2920): persist positions on TagBegin, Class_, TagEnd. */
		openTagBegin := stdbytes.LastIndexByte(h.Response.Body, '<')
		if openTagBegin == -1 {
			panic("failed to find open angle bracket")
		}
		if h.Response.Body[openTagBegin+1] == '/' {
			closeTagBegin = openTagBegin
			openTagBegin = stdbytes.LastIndexByte(h.Response.Body[:closeTagBegin], '<')
			if openTagBegin == -1 {
				panic("failed to find open angle bracket")
			}
		}

		openTagEnd := stdbytes.IndexByte(h.Response.Body[openTagBegin:], '>')
		if openTagEnd == -1 {
			panic("failed to find close angle bracket")
		}
		openTagEnd += openTagBegin

		classPos := stdbytes.Index(h.Response.Body[openTagBegin:openTagEnd], _class)
		if (classPos == -1) && (closeTagBegin == -1) {
			h.Backspace().String(` class="`).String(class).String(`">`)
		} else if classPos == -1 {
			backup := h.Response.Arena.Copy(h.Response.Body[openTagEnd:])
			h.Response.Body = h.Response.Body[:openTagEnd]
			h.String(` class="`).String(class).String(`"`).Bytes(backup)
		} else {
			classPos += openTagBegin + len(`class="`)
			quote := stdbytes.IndexByte(h.Response.Body[classPos:], '"')
			if quote == -1 {
				panic("failed to find end quote in class definitions")
			}
			quote += classPos

			backup := h.Response.Arena.Copy(h.Response.Body[quote:])
			h.Response.Body = h.Response.Body[:quote]
			h.String(" ").String(class).Bytes(backup)
		}
	}
	return h
}

func (h *HTML) Classes(classes ...string) *HTML {
	var required int
	for i := 0; i < len(classes); i++ {
		required += len(classes[i])
	}

	if required > 0 {
		var n int

		buf := h.Response.Arena.NewSlice(required + len(classes))
		for i := 0; i < len(classes); i++ {
			n += copy(buf[n:], classes[i])
			buf[n] = ' '
			n++
		}
		n--

		h.Class_(bytes.AsString(buf[:n]))
	}

	return h
}

func (h *HTML) Action(action string) *HTML {
	return h.Backspace().String(` action="`).String(action).String(`">`)
}

func (h *HTML) AriaCurrent(current string) *HTML {
	return h.Backspace().String(` aria-current="`).String(current).String(`">`)
}

func (h *HTML) AriaLabel(label string) *HTML {
	return h.Backspace().String(` aria-label="`).String(label).String(`">`)
}

func (h *HTML) Checked(checked bool) *HTML {
	if checked {
		h.Backspace().String(" checked>")
	}
	return h
}

func (h *HTML) DataBsToggle(toggle string) *HTML {
	return h.Backspace().String(` data-bs-toggle="`).String(toggle).String(`">`)
}

func (h *HTML) DominantBaseline(baseline string) *HTML {
	return h.Backspace().String(` dominant-baseline="`).String(baseline).String(`">`)
}

func (h *HTML) Fill(fill string) *HTML {
	return h.Backspace().String(` fill="`).String(fill).String(`">`)
}

func (h *HTML) FontFamily(family string) *HTML {
	return h.Backspace().String(` font-family="`).String(family).String(`">`)
}

func (h *HTML) FontSize(size int) *HTML {
	return h.Backspace().String(` font-size="`).Int(size).String(`">`)
}

func (h *HTML) FontWeight(weight int) *HTML {
	return h.Backspace().String(` font-weight="`).Int(weight).String(`">`)
}

func (h *HTML) For(for_ string) *HTML {
	return h.Backspace().String(` for="`).String(for_).String(`">`)
}

func (h *HTML) Id(id string) *HTML {
	return h.Backspace().String(` id="`).String(id).String(`">`)
}

func (h *HTML) Href(href string) *HTML {
	return h.Backspace().String(` href="`).String(href).String(`">`)
}

func (h *HTML) Minlength(minLength int) *HTML {
	return h.Backspace().String(` minlength="`).Int(minLength).String(`">`)
}

func (h *HTML) Maxlength(maxLength int) *HTML {
	return h.Backspace().String(` maxlength="`).Int(maxLength).String(`">`)
}

func (h *HTML) Name(name string) *HTML {
	return h.Backspace().String(` name="`).String(name).String(`">`)
}

func (h *HTML) Onclick(onclick string) *HTML {
	return h.Backspace().String(` onclick="`).String(onclick).String(`">`)
}

func (h *HTML) Placeholder(placeholder string) *HTML {
	return h.Backspace().String(` placeholder="`).LString(placeholder).String(`">`)
}

func (h *HTML) Rel(rel string) *HTML {
	return h.Backspace().String(` rel="`).String(rel).String(`">`)
}

func (h *HTML) Required(required bool) *HTML {
	if required {
		h.Backspace().String(" required>")
	}
	return h
}

func (h *HTML) Selected(selected bool) *HTML {
	if selected {
		h.Backspace().String(" selected>")
	}
	return h
}

func (h *HTML) Stroke(stroke string) *HTML {
	return h.Backspace().String(` stroke="`).String(stroke).String(`">`)
}

func (h *HTML) StrokeLinecap(linecap string) *HTML {
	return h.Backspace().String(` stroke-linecap="`).String(linecap).String(`">`)
}

func (h *HTML) StrokeLinejoin(linejoin string) *HTML {
	return h.Backspace().String(` stroke-linejoin="`).String(linejoin).String(`">`)
}

func (h *HTML) StrokeWidth(width int) *HTML {
	return h.Backspace().String(` stroke-width="`).Int(width).String(`">`)
}

func (h *HTML) Value(value string) *HTML {
	return h.Backspace().String(` value="`).LString(value).String(`">`)
}

func (h *HTML) Viewbox(viewbox string) *HTML {
	return h.Backspace().String(` viewbox="`).String(viewbox).String(`">`)
}
