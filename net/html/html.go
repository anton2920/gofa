package html

import (
	"github.com/anton2920/gofa/l10n"
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/time"
	"github.com/anton2920/gofa/trace"
)

type Theme struct {
	Form   Attributes
	Label  Attributes
	Input  Attributes
	Button Attributes
}

type HTML struct {
	Theme

	*http.Response
	l10n.Language
	time.Timezone
}

func New(w *http.Response, l l10n.Language, tz time.Timezone) HTML {
	var h HTML

	h.Response = w
	h.Language = l
	h.Timezone = tz

	return h
}

func (h *HTML) Int(n int) {
	h.WriteInt(n)
}

func (h *HTML) HTMLString(s string) {
	h.WriteHTMLString(s)
}

func (h *HTML) LString(s string) {
	h.WriteHTMLString(h.L(s))
}

func (h *HTML) String(s string) (int, error) {
	return h.WriteString(s)
}

func (h *HTML) TagBegin(tag string, attrs ...Attributes) {
	t := trace.Begin("")

	attr := h.MergeAttributes(attrs...)

	h.String(` <`)
	h.String(tag)

	if len(attrs) > 0 {
		DisplayStringAttribute(h, "class", attr.Class)

		DisplayStringAttribute(h, "action", attr.Action)
		DisplayStringAttribute(h, "enctype", attr.Enctype)
		DisplayStringAttribute(h, "href", attr.Href)
		DisplayStringAttribute(h, "id", attr.ID)
		DisplayStringAttribute(h, "method", attr.Method)
		DisplayStringAttribute(h, "name", attr.Name)
		DisplayStringAttribute(h, "src", attr.Src)
		DisplayStringAttribute(h, "type", attr.Type)

		DisplayLStringAttribute(h, "alt", attr.Alt)
		DisplayLStringAttribute(h, "placeholder", attr.Placeholder)
		DisplayLStringAttribute(h, "value", attr.Value)

		DisplayIntAttribute(h, "cols", attr.Cols)
		DisplayIntAttribute(h, "max", attr.Max)
		DisplayIntAttribute(h, "maxlength", attr.MaxLength)
		DisplayIntAttribute(h, "min", attr.Min)
		DisplayIntAttribute(h, "minlength", attr.MinLength)
		DisplayIntAttribute(h, "rows", attr.Rows)

		DisplayBoolAttribute(h, "disabled", attr.Disabled)
		DisplayBoolAttribute(h, "formnovalidate", attr.FormNoValidate)
		DisplayBoolAttribute(h, "readonly", attr.Readonly)
		DisplayBoolAttribute(h, "required", attr.Required)
	}
	h.String(`> `)

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
	h.String(`">`)
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
	h.TagBegin("body", attrs...)
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

func (h *HTML) DivBegin(class string) {
	h.TagBegin("div", Attributes{Class: class})
}

func (h *HTML) DivEnd() {
	h.TagEnd("div")
}

func (h *HTML) Img(alt string, src string, attrs ...Attributes) {
	h.TagBegin("img", h.AppendAttributes(attrs, Attributes{Alt: alt, Src: src}))
}

func (h *HTML) ABegin(href string, attrs ...Attributes) {
	h.TagBegin("a", h.AppendAttributes(attrs, Attributes{Href: href}))
}

func (h *HTML) A(href string, value string, attrs ...Attributes) {
	h.ABegin(href, attrs...)
	h.LString(value)
	h.AEnd()
}

func (h *HTML) AEnd() {
	h.TagEnd("a")
}
