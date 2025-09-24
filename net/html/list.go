package html

func (h *HTML) OLBegin(attrs ...Attributes) {
	h.TagBegin("ol", h.PrependAttributes(h.Theme.OL, attrs))
}

func (h *HTML) OLEnd() {
	h.TagEnd("ol")
}

func (h *HTML) ULBegin(attrs ...Attributes) {
	h.TagBegin("ul", h.PrependAttributes(h.Theme.UL, attrs))
}

func (h *HTML) ULEnd() {
	h.TagEnd("ul")
}

func (h *HTML) LIBegin(attrs ...Attributes) {
	h.TagBegin("li", h.PrependAttributes(h.Theme.LI, attrs))
}

func (h *HTML) LI(li string, attrs ...Attributes) {
	h.LIBegin(attrs...)
	h.LString(li)
	h.LIEnd()
}

func (h *HTML) LIEnd() {
	h.TagEnd("li")
}
