package html

func (h *HTML) H1Begin(attrs ...Attributes) {
	h.TagBegin("h1", attrs...)
}

func (h *HTML) H1(s string, attrs ...Attributes) {
	h.H1Begin(attrs...)
	h.HTMLString(s)
	h.H1End()
}

func (h *HTML) H1End() {
	h.TagEnd("h1")
}

func (h *HTML) H2Begin(attrs ...Attributes) {
	h.TagBegin("h2", attrs...)
}

func (h *HTML) H2(s string, attrs ...Attributes) {
	h.H2Begin(attrs...)
	h.HTMLString(s)
	h.H2End()
}

func (h *HTML) H2End() {
	h.TagEnd("h2")
}

func (h *HTML) H3Begin(attrs ...Attributes) {
	h.TagBegin("h3", attrs...)
}

func (h *HTML) H3(s string, attrs ...Attributes) {
	h.H3Begin(attrs...)
	h.HTMLString(s)
	h.H3End()
}

func (h *HTML) H3End() {
	h.TagEnd("h3")
}

func (h *HTML) H4Begin(attrs ...Attributes) {
	h.TagBegin("h4", attrs...)
}

func (h *HTML) H4(s string, attrs ...Attributes) {
	h.H4Begin(attrs...)
	h.HTMLString(s)
	h.H4End()
}

func (h *HTML) H4End() {
	h.TagEnd("h4")
}

func (h *HTML) H5Begin(attrs ...Attributes) {
	h.TagBegin("h5", attrs...)
}

func (h *HTML) H5(s string, attrs ...Attributes) {
	h.H5Begin(attrs...)
	h.HTMLString(s)
	h.H5End()
}

func (h *HTML) H5End() {
	h.TagEnd("h5")
}

func (h *HTML) H6Begin(attrs ...Attributes) {
	h.TagBegin("h6", attrs...)
}

func (h *HTML) H6(s string, attrs ...Attributes) {
	h.H6Begin(attrs...)
	h.HTMLString(s)
	h.H6End()
}

func (h *HTML) H6End() {
	h.TagEnd("h6")
}
