package html

import (
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/l10n"
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/time"
	"github.com/anton2920/gofa/trace"
)

type HTML struct {
	W *http.Response

	database.ID
	l10n.Language
	time.Timezone
}

func New(w *http.Response, id database.ID, l l10n.Language, tz time.Timezone) HTML {
	var h HTML

	h.ID = id
	h.Language = l
	h.Timezone = tz

	return h
}

func (h *HTML) Int(n int) {
	h.W.WriteInt(n)
}

func (h *HTML) HTMLString(s string) {
	h.W.WriteHTMLString(s)
}

func (h *HTML) LocalizedString(s string) (int, error) {
	return h.W.WriteString(h.L(s))
}

func (h *HTML) String(s string) (int, error) {
	return h.W.WriteString(s)
}

func (h *HTML) TagBegin(tag string, attrs ...Attributes) {
	t := trace.Begin("")

	attr := h.MergeAttributes(attrs...)

	h.String(` <`)
	h.String(tag)

	if len(attrs) > 0 {
		DisplayStringAttribute(h, "class", attr.Class)

		DisplayStringAttribute(h, "alt", attr.Alt)
		DisplayStringAttribute(h, "src", attr.Src)
		DisplayStringAttribute(h, "id", attr.ID)
		DisplayStringAttribute(h, "name", attr.Name)
		DisplayStringAttribute(h, "placeholder", attr.Placeholder)
		DisplayStringAttribute(h, "type", attr.Type)
		DisplayStringAttribute(h, "value", attr.Value)
		DisplayStringAttribute(h, "method", attr.Method)
		DisplayStringAttribute(h, "action", attr.Action)
		DisplayStringAttribute(h, "enctype", attr.Enctype)

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
	h.HTMLString(h.L(title))
	h.TitleEnd()
}

func (h *HTML) TitleEnd() {
	h.TagEnd("title")
}

func (h *HTML) BodyBegin() {
	h.TagBegin("body")
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

func (h *HTML) DivBegin(attrs ...Attributes) {
	h.TagBegin("div", attrs...)
}

func (h *HTML) DivEnd() {
	h.TagEnd("div")
}

func (h *HTML) Img(alt string, src string, attrs ...Attributes) {
	h.TagBegin("img", h.AppendAttributes(attrs, Attributes{Alt: alt, Src: src}))
}
