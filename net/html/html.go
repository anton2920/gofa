package html

import (
	stdtime "time"
	"unicode"

	"github.com/anton2920/gofa/bools"
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/debug"
	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/l10n"
	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/session"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/time"
	"github.com/anton2920/gofa/trace"
)

type Theme struct {
	A        Attributes
	Body     Attributes
	Button   Attributes
	Checkbox Attributes
	Div      Attributes
	Form     Attributes
	H1       Attributes
	H2       Attributes
	H3       Attributes
	H4       Attributes
	H5       Attributes
	H6       Attributes
	Img      Attributes
	Input    Attributes
	LI       Attributes
	Label    Attributes
	OL       Attributes
	P        Attributes
	Select   Attributes
	Span     Attributes
	Textarea Attributes
	UL       Attributes

	Pagination             Attributes
	PaginationButton       Attributes
	PaginationButtonActive Attributes
}

type HTML struct {
	*http.Response
	*http.Request

	session.Session
	Theme Theme
}

func New(w *http.Response, r *http.Request, session session.Session, theme Theme) HTML {
	return HTML{Response: w, Request: r, Session: session, Theme: theme}
}

/* TODO(anton2920): check whether it allocates memory. */
func (h *HTML) WithoutTheme() *HTML {
	return &HTML{Response: h.Response, Request: h.Request, Session: h.Session}
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
	h.LString(s[:len(s)-bools.ToInt(n == 1)])
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
	buf := h.Arena.NewSlice(len(stdtime.DateOnly))
	stdtime.Unix(d+int64(h.Timezone)*time.OneHour, 0).UTC().AppendFormat(buf[:0], stdtime.DateOnly)
	return bytes.AsString(buf)
}

func (h *HTML) Itoa(x int) string {
	buf := h.Arena.NewSlice(ints.Bufsize)
	n := slices.PutInt(buf, x)
	return bytes.AsString(buf[:n])
}

func (h *HTML) Itoa1(x int) string {
	if x == 0 {
		return ""
	}
	return h.Itoa(x)
}

func (h *HTML) Ttoa(d int64) string {
	buf := h.Arena.NewSlice(len(stdtime.DateTime))
	stdtime.Unix(d+int64(h.Timezone)*time.OneHour, 0).UTC().AppendFormat(buf[:0], stdtime.DateTime)
	return bytes.AsString(buf)
}

func (h *HTML) IndexedName(name string, index int) string {
	var n int

	buf := h.Arena.NewSlice(len(name) + ints.Bufsize)
	n += copy(buf[n:], name)
	n += slices.PutInt(buf[n:], index)

	return bytes.AsString(buf[:n])
}

func (h *HTML) DoublyIndexedName(name string, index1 int, index2 int) string {
	var n int

	buf := h.Arena.NewSlice(len(name) + ints.Bufsize + ints.Bufsize + 1)
	n += copy(buf[n:], name)
	n += slices.PutInt(buf[n:], index1)
	n += copy(buf[n:], ".")
	n += slices.PutInt(buf[n:], index2)

	return bytes.AsString(buf[:n])
}

func (h *HTML) PathWithID(path string, id database.ID) string {
	var n int

	buf := h.Arena.NewSlice(len(path) + ints.Bufsize)
	n += copy(buf[n:], path)
	n += slices.PutInt(buf[n:], int(id))

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
	h.String(`</`)
	h.String(tag)
	h.String(`>`)
}

func (h *HTML) Begin() {
	h.String(`<!DOCTYPE html>`)
	h.String(`<html lang="`)
	h.String(l10n.Language2HTMLLang[h.Language])
	h.String(`" data-bs-theme="light">`)
}

func (h *HTML) End() {
	h.TagEnd("html")
}

func (h *HTML) HeadBegin() {
	h.TagBegin("head")
	h.String(`<meta charset="UTF-8">`)
	h.String(`<meta name="viewport" content="width=device-width, initial-scale=1.0">`)
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
	h.ABegin(href, attrs...)
	h.LString(value)
	h.AEnd()
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
			h.StatusCode = httpError.StatusCode
			message = httpError.DisplayErrorMessage
		} else if _, ok := err.(errors.Panic); ok {
			h.StatusCode = http.StatusInternalServerError
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
		h.DivBegin(attrs...)
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
