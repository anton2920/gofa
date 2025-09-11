package html

import (
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/slices"
)

func (h *HTML) FormBegin(method string, action string, attrs ...Attributes) {
	h.TagBegin("form", h.AppendAttributes(attrs, Attributes{Method: method, Action: action}))
}

func (h *HTML) FormEnd() {
	h.TagEnd("form")
}

func (h *HTML) Input(typ string, attrs ...Attributes) {
	h.TagBegin("input", h.AppendAttributes(attrs, Attributes{Type: typ}))
}

func (h *HTML) HiddenBool(name string, b bool) {
	if b {
		h.Input("hidden", Attributes{Name: name})
	}
}

func (h *HTML) HiddenID(name string, id database.ID) {
	h.HiddenInt(name, int(id))
}

func (h *HTML) HiddenInt(name string, n int) {
	if n > 0 {
		buffer := make([]byte, ints.Bufsize)
		slices.PutInt(buffer, n)
		h.Input("hidden", Attributes{Name: bytes.AsString(buffer)})
	}
}

func (h *HTML) HiddenString(name string, s string) {
	if len(s) > 0 {
		h.Input("hidden", Attributes{Name: name, Value: s})
	}
}
