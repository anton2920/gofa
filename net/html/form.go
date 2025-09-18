package html

import "github.com/anton2920/gofa/database"

func (h *HTML) FormBegin(method string, action string, attrs ...Attributes) {
	h.TagBegin("form", h.Theme.Form, h.AppendAttributes(attrs, Attributes{Method: method, Action: action}))
}

func (h *HTML) FormEnd() {
	h.TagEnd("form")
}

func (h *HTML) Input(typ string, attrs ...Attributes) {
	h.TagBegin("input", h.Theme.Input, h.AppendAttributes(attrs, Attributes{Type: typ}))
}

func (h *HTML) HiddenBool(name string, b bool) {
	if b {
		h.Input("hidden", Attributes{Name: name})
	}
}

func (h *HTML) HiddenID(name string, id database.ID) {
	h.HiddenInt(name, int(id))
}

func (h *HTML) HiddenInt(name string, x int) {
	if x > 0 {
		h.Input("hidden", Attributes{Name: name, Value: h.Itoa(x)})
	}
}

func (h *HTML) HiddenString(name string, s string) {
	if len(s) > 0 {
		h.Input("hidden", Attributes{Name: name, Value: s})
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

func (h *HTML) SelectBegin(attrs ...Attributes) {
	h.TagBegin("select", h.PrependAttributes(h.Theme.Select, attrs))
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
