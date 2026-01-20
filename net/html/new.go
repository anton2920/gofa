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
	return h.String(` <head>`)
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

func (h *HTML) ABegin2(href string) *HTML {
	return h.String(` <a href="`).String(href).String(`">`).Class_(h.Theme.A.Class)
}

func (h *HTML) A2(href string, contents string) *HTML {
	return h.ABegin2(href).LString(contents).AEnd2()
}

func (h *HTML) AEnd2() *HTML {
	return h.String(`</a>`)
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

func (h *HTML) Checked(checked bool) *HTML {
	if checked {
		h.Backspace().String(" checked>")
	}
	return h
}

func (h *HTML) HREF(href string) *HTML {
	return h.Backspace().String(` href="`).String(href).String(`">`)
}

func (h *HTML) MinLength(minLength int) *HTML {
	return h.Backspace().String(` minlength="`).Int(minLength).String(`">`)
}

func (h *HTML) MaxLength(maxLength int) *HTML {
	return h.Backspace().String(` maxlength="`).Int(maxLength).String(`">`)
}

func (h *HTML) Name(name string) *HTML {
	return h.Backspace().String(` name="`).String(name).String(`">`)
}

func (h *HTML) Placeholder(placeholder string) *HTML {
	return h.Backspace().String(` placeholder="`).String(placeholder).String(`">`)
}

func (h *HTML) Required(required bool) *HTML {
	if required {
		h.Backspace().String(" required>")
	}
	return h
}

func (h *HTML) Value(value string) *HTML {
	return h.Backspace().String(` value="`).String(value).String(`">`)
}
