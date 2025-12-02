package html

import (
	stdtime "time"
	"unicode"

	"github.com/anton2920/gofa/bools"
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/debug"
	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/l10n"
	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/session"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/time"
	"github.com/anton2920/gofa/trace"
)

type HTML struct {
	*http.Response
	*http.Request

	withoutTheme *HTML

	Theme Theme
}

func New(w *http.Response, r *http.Request, theme Theme) HTML {
	return HTML{Response: w, Request: r, Theme: theme}
}

func (h *HTML) WithoutTheme() *HTML {
	if h.withoutTheme == nil {
		h.withoutTheme = new(HTML)
		h.withoutTheme.Response = h.Response
		h.withoutTheme.Request = h.Request
	}
	return h.withoutTheme
}

func (h *HTML) Bytes(bs []byte) {
	h.Write(bs)
}

func (h *HTML) Date(d int64) {
	h.String(h.Dtoa(d))
}

func (h *HTML) Int(n int) {
	h.WriteInt(n)
}

func (h *HTML) HString(s string) {
	h.WriteHTMLString(s)
}

func (h *HTML) LString(s string) {
	h.HString(h.L(s))
}

func (h *HTML) LStringColon(s string) {
	h.LString(s)
	h.String(": ")
}

func (h *HTML) LStringPlural(s string, n int) {
	const suffix = "ies"

	if !strings.EndsWith(s, suffix) {
		h.LString(s[:len(s)-bools.ToInt(n == 1)])
	} else {
		buf := h.Response.Arena.NewSlice(len(s))
		copy(buf, s)

		if n == 1 {
			buf = buf[:len(s)-len(suffix)+1]
			buf[len(buf)-1] = 'y'
		}

		h.LString(bytes.AsString(buf))
	}
}

func (h *HTML) String(s string) {
	h.WriteString(s)
}

func (h *HTML) Time(t int64) {
	h.String(h.Ttoa(t))
}

func (h *HTML) TString(s string) {
	if s := h.L(s); len(s) > 1 {
		h.String(string(unicode.ToUpper(rune(s[0]))))
		h.String(s[1:])
	}
}

func (h *HTML) Dtoa(d int64) string {
	const format = "2006-01-02"
	v := d + int64(h.Timezone)*time.Hour
	return stdtime.Unix(v/time.Second, v%time.Second).UTC().Format(format)
}

func (h *HTML) Dtoa1(d int64) string {
	if d == 0 {
		return ""
	}
	return h.Dtoa(d)
}

func (h *HTML) Itoa(x int) string {
	buf := h.Response.Arena.NewSlice(ints.Bufsize)
	n := slices.PutInt(buf, x)
	return bytes.AsString(buf[:n])
}

func (h *HTML) Itoa1(x int) string {
	if x == 0 {
		return ""
	}
	return h.Itoa(x)
}

/* TODO(anton2920): unify this and 'Dtoa'. */
func (h *HTML) Ttoa(t int64) string {
	const format = "2006-01-02 15:04:05"
	v := t + int64(h.Timezone)*time.Hour
	return stdtime.Unix(v/time.Second, v%time.Second).UTC().Format(format)
}

func (h *HTML) IndexedName(name string, indicies ...int) string {
	var n int

	buf := h.Response.Arena.NewSlice(len(name) + ints.Bufsize*len(indicies) + 2*len(indicies))
	n += copy(buf[n:], name)
	for i := 0; i < len(indicies); i++ {
		index := indicies[i]
		if i > 0 {
			n += copy(buf[n:], ".")
		}
		n += slices.PutInt(buf[n:], index)
	}

	return bytes.AsString(buf[:n])
}

func (h *HTML) TagBegin(tag string, attrs ...Attributes) {
	t := trace.Begin("")

	attr := h.MergeAttributes(attrs...)

	h.String(` <`)
	h.String(tag)

	if len(attrs) > 0 {
		DisplayStringAttribute(h, "class", attr.Class)

		DisplayStringAttribute(h, "accept", attr.Accept)
		DisplayStringAttribute(h, "action", attr.Action)
		DisplayStringAttribute(h, "enctype", attr.Enctype)
		DisplayStringAttribute(h, "href", attr.Href)
		DisplayStringAttribute(h, "id", attr.ID)
		DisplayStringAttribute(h, "method", attr.Method)
		DisplayStringAttribute(h, "name", attr.Name)
		DisplayStringAttribute(h, "onclick", attr.OnClick)
		DisplayStringAttribute(h, "rel", attr.Rel)
		DisplayStringAttribute(h, "src", attr.Src)
		DisplayStringAttribute(h, "style", attr.Style)
		DisplayStringAttribute(h, "type", attr.Type)
		DisplayStringAttribute(h, "value", attr.Value)

		DisplayLStringAttribute(h, "alt", attr.Alt)
		DisplayLStringAttribute(h, "placeholder", attr.Placeholder)

		DisplayIntAttribute(h, "cols", attr.Cols)
		DisplayIntAttribute(h, "max", attr.Max)
		DisplayIntAttribute(h, "maxlength", attr.MaxLength)
		DisplayIntAttribute(h, "min", attr.Min)
		DisplayIntAttribute(h, "minlength", attr.MinLength)
		DisplayIntAttribute(h, "rows", attr.Rows)

		DisplayBoolAttribute(h, "checked", attr.Checked)
		DisplayBoolAttribute(h, "disabled", attr.Disabled)
		DisplayBoolAttribute(h, "formnovalidate", attr.FormNoValidate)
		DisplayBoolAttribute(h, "readonly", attr.Readonly)
		DisplayBoolAttribute(h, "required", attr.Required)
		DisplayBoolAttribute(h, "selected", attr.Selected)
	}
	h.String(`>`)

	trace.End(t)
}

func (h *HTML) TagEnd(tag string) {
	t := trace.Begin("")

	h.String(`</`)
	h.String(tag)
	h.String(`>`)

	trace.End(t)
}

func (h *HTML) Begin() {
	h.String(`<!DOCTYPE html>`)
	h.String(`<html lang="`)
	h.String(l10n.Language2HTMLLang[h.Language])
	h.String(`"`)
	if h.ColorScheme > 0 {
		h.String(` data-bs-theme="`)
		h.String(session.ColorScheme2String[h.ColorScheme])
		h.String(`"`)
	}
	h.String(`>`)
}

func (h *HTML) End() {
	h.TagEnd("html")
}

func (h *HTML) HeadBegin() {
	h.TagBegin("head")
	h.String(`<meta charset="UTF-8">`)
	h.String(`<meta name="viewport" content="width=device-width, initial-scale=1.0">`)

	h.Link(h.Theme.HeadLink)
	h.Script("", h.Theme.HeadScript)
}

func (h *HTML) HeadEnd() {
	h.TagEnd("head")
}

func (h *HTML) TitleBegin() {
	h.TagBegin("title")
}

func (h *HTML) Title(title string) {
	h.TitleBegin()
	h.LString(title)
	h.TitleEnd()
}

func (h *HTML) TitleEnd() {
	h.TagEnd("title")
}

func (h *HTML) Link(attrs ...Attributes) {
	h.TagBegin("link", attrs...)
}

func (h *HTML) ScriptBegin(attrs ...Attributes) {
	h.TagBegin("script", attrs...)
}

func (h *HTML) Script(script string, attrs ...Attributes) {
	h.ScriptBegin(attrs...)
	h.String(script)
	h.ScriptEnd()
}

func (h *HTML) ScriptEnd(attrs ...Attributes) {
	h.TagEnd("script")
}

func (h *HTML) BodyBegin(attrs ...Attributes) {
	h.TagBegin("body", h.PrependAttributes(h.Theme.Body, attrs))
}

func (h *HTML) BodyEnd() {
	h.TagEnd("body")
}

func (h *HTML) BR() {
	h.String(`<br>`)
}

func (h *HTML) HR() {
	h.String(`<hr>`)
}

func (h *HTML) SP() {
	h.String(` `)
}

func (h *HTML) ABegin(href string, attrs ...Attributes) {
	h.TagBegin("a", h.Theme.A, h.AppendAttributes(attrs, Attributes{Href: href}))
}

func (h *HTML) A(href string, value string, attrs ...Attributes) {
	if len(href) > 0 {
		h.ABegin(href, attrs...)
		h.LString(value)
		h.AEnd()
	}
}

func (h *HTML) AEnd() {
	h.TagEnd("a")
}

func (h *HTML) DivBegin(attrs ...Attributes) {
	h.TagBegin("div", h.PrependAttributes(h.Theme.Div, attrs))
}

func (h *HTML) DivEnd() {
	h.TagEnd("div")
}

func (h *HTML) Error(err error, attrs ...Attributes) {
	var message string

	if err != nil {
		if httpError, ok := err.(http.Error); ok {
			h.Status = httpError.Status
			message = httpError.DisplayErrorMessage
		} else if _, ok := err.(errors.Panic); ok {
			h.Status = http.StatusInternalServerError
			message = http.ServerDisplayErrorMessage
		} else {
			log.Panicf("Unsupported error type: %T (%v)", err, err)
		}

		if debug.Debug {
			message = err.Error()
		}
	}

	h.ErrorMessage(message, attrs...)
}

func (h *HTML) ErrorMessage(message string, attrs ...Attributes) {
	if len(message) > 0 {
		h.DivBegin(h.PrependAttributes(h.Theme.Error, attrs))
		h.LStringColon("Error")
		h.LString(message)
		h.DivEnd()
	}
}

func (h *HTML) Img(alt string, src string, attrs ...Attributes) {
	h.TagBegin("img", h.Theme.Img, h.AppendAttributes(attrs, Attributes{Alt: alt, Src: src}))
}

func (h *HTML) PBegin(attrs ...Attributes) {
	h.TagBegin("p", h.PrependAttributes(h.Theme.P, attrs))
}

func (h *HTML) P(p string, attrs ...Attributes) {
	if len(p) > 0 {
		h.PBegin(attrs...)
		h.LString(p)
		h.PEnd()
	}
}

func (h *HTML) PEnd() {
	h.TagEnd("p")
}

func (h *HTML) SpanBegin(attrs ...Attributes) {
	h.TagBegin("span", h.PrependAttributes(h.Theme.Span, attrs))
}

func (h *HTML) Span(s string, attrs ...Attributes) {
	if len(s) > 0 {
		h.SpanBegin(attrs...)
		h.LString(s)
		h.SpanEnd()
	}
}

func (h *HTML) SpanEnd() {
	h.TagEnd("span")
}
