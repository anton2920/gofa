package html

import (
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/time"
)

func (h *HTML) FormBegin(method string, attrs ...Attributes) {
	h.TagBegin("form", h.Theme.Form, h.AppendAttributes(attrs, Attributes{Method: method}))
}

func (h *HTML) FormEnd() {
	h.TagEnd("form")
}

func (h *HTML) Input(typ string, attrs ...Attributes) {
	h.TagBegin("input", h.Theme.Input, h.AppendAttributes(attrs, Attributes{Type: typ}))
}

func (h *HTML) HiddenBool(name string, b bool) {
	if b {
		h.WithoutTheme().Input("hidden", Attributes{Name: name})
	}
}

func (h *HTML) HiddenID(name string, id database.ID) {
	h.HiddenInt(name, int(id))
}

func (h *HTML) HiddenInt(name string, x int) {
	if x > 0 {
		h.WithoutTheme().Input("hidden", Attributes{Name: name, Value: h.Itoa(x)})
	}
}

func (h *HTML) HiddenString(name string, s string) {
	if len(s) > 0 {
		h.WithoutTheme().Input("hidden", Attributes{Name: name, Value: s})
	}
}

func (h *HTML) LabelBegin(attrs ...Attributes) {
	h.TagBegin("label", h.PrependAttributes(h.Theme.Label, attrs))
}

func (h *HTML) Label(text string, attrs ...Attributes) {
	h.LabelBegin(attrs...)
	h.LString(text)
	h.LabelEnd()
}

func (h *HTML) LabelEnd() {
	h.TagEnd("label")
}

func (h *HTML) Button(value string, attrs ...Attributes) {
	h.WithoutTheme().Input("submit", h.Theme.Button, h.AppendAttributes(attrs, Attributes{Value: h.L(value)}))
}

func (h *HTML) Checkbox(attrs ...Attributes) {
	h.WithoutTheme().Input("checkbox", h.PrependAttributes(h.Theme.Checkbox, attrs))
}

func (h *HTML) OptionBegin(attrs ...Attributes) {
	h.TagBegin("option", attrs...)
}

func (h *HTML) Option(opt string, attrs ...Attributes) {
	h.OptionBegin(attrs...)
	h.LString(opt)
	h.OptionEnd()
}

func (h *HTML) OptionEnd() {
	h.TagEnd("option")
}

func (h *HTML) SelectBegin(attrs ...Attributes) {
	h.TagBegin("select", h.PrependAttributes(h.Theme.Select, attrs))
}

func (h *HTML) Select(xs []string, start int, selected int, attrs ...Attributes) {
	res := h.MergeAttributes(attrs...)

	placeholder := res.Placeholder
	res.Placeholder = ""

	required := res.Required
	res.Required = false

	h.SelectBegin(res)
	if !required {
		h.Option(placeholder)
	}
	for i := start; i < len(xs); i++ {
		h.Option(xs[i], Attributes{Value: h.Itoa(i), Selected: i == selected})
	}
	h.SelectEnd()
}

func (h *HTML) Select0(xs []string, selected int, attrs ...Attributes) {
	h.Select(xs, 0, selected, h.AppendAttributes(attrs, Attributes{Required: true}))
}

func (h *HTML) Select1(xs []string, selected int, attrs ...Attributes) {
	h.Select(xs, 1, selected, attrs...)
}

func (h *HTML) TimezoneSelect(selected time.Timezone, attrs ...Attributes) {
	h.SelectBegin(attrs...)
	for i := time.TimezoneNone + 1; i < time.TimezoneCount; i++ {
		h.Option(time.Timezone2String[i], Attributes{Value: h.Itoa(int(i)), Selected: i == selected})
	}
	h.SelectEnd()
}

func (h *HTML) SelectEnd(attrs ...Attributes) {
	h.TagEnd("select")
}

func (h *HTML) TextareaBegin(attrs ...Attributes) {
	h.TagBegin("textarea", h.PrependAttributes(h.Theme.Textarea, attrs))
}

func (h *HTML) Textarea(value string, attrs ...Attributes) {
	h.TextareaBegin(attrs...)
	h.HString(value)
	h.TextareaEnd()
}

func (h *HTML) TextareaEnd() {
	h.TagEnd("textarea")
}
