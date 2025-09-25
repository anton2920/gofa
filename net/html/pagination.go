package html

import "github.com/anton2920/gofa/bools"

type Pagination struct {
	CurrentPage   int
	ItemsPerPage  int
	NumberOfPages int
	WindowSize    int
}

const (
	PaginationKeyCurrentPage  = "Page"
	PaginationKeyItemsPerPage = "ItemsPerPage"
)

func (pagination *Pagination) GetItemsStartPosition() int {
	return pagination.ItemsPerPage * (pagination.CurrentPage - bools.ToInt(pagination.CurrentPage > 0))
}

func (pagination *Pagination) UpdateNumberOfPages(nitems int) {
	pagination.NumberOfPages = (nitems / pagination.ItemsPerPage) + bools.ToInt(nitems%pagination.ItemsPerPage > 0)
}

func (h *HTML) PageSelectorButton(text string, page int, active bool) {
	attrs := h.Theme.PageSelectorButtonContainer
	if active {
		attrs = h.MergeAttributes(attrs, h.Theme.PageSelectorButtonActive)
	}

	h.LIBegin(attrs)
	if text == "" {
		// h.A("/", h.Itoa(page))
		h.WithoutTheme().Button(h.Itoa(page), h.Theme.PageSelectorButton, Attributes{Name: PaginationKeyCurrentPage, FormNoValidate: true})
	} else {
		h.TagBegin("button", h.Theme.PageSelectorButton, Attributes{Name: PaginationKeyCurrentPage, Value: h.Itoa(page), FormNoValidate: true})
		h.String(text)
		h.TagEnd("button")
	}
	h.LIEnd()
}

func (h *HTML) PageSelectorEllipsis() {
	h.WithoutTheme().Button("...", h.Theme.PageSelectorButton, Attributes{Disabled: true})
}

func (h *HTML) PageSelector(pagination *Pagination, attrs ...Attributes) {
	if pagination.NumberOfPages > 1 {
		if pagination.CurrentPage == 0 {
			pagination.CurrentPage = 1
		}
		h.DivBegin(h.PrependAttributes(h.Theme.PageSelector, attrs))

		prev := pagination.CurrentPage - bools.ToInt(pagination.CurrentPage > 1)
		h.PageSelectorButton("&laquo;", prev, false)

		if pagination.WindowSize > pagination.NumberOfPages {
			pagination.WindowSize = pagination.NumberOfPages - 1
		}

		if pagination.NumberOfPages == 2 {
			for i := 1; i <= pagination.NumberOfPages; i++ {
				h.PageSelectorButton("", i, i == pagination.CurrentPage)
			}
		} else {
			switch {
			case pagination.CurrentPage < pagination.WindowSize:
				for i := 1; i <= pagination.WindowSize; i++ {
					h.PageSelectorButton("", i, i == pagination.CurrentPage)
				}
				h.PageSelectorEllipsis()
				h.PageSelectorButton("", pagination.NumberOfPages, pagination.NumberOfPages == pagination.CurrentPage)
			default:
				h.PageSelectorButton("", 1, pagination.CurrentPage == 1)
				h.PageSelectorEllipsis()

				halfWindow := pagination.WindowSize / 2
				start := pagination.CurrentPage - halfWindow + bools.ToInt((pagination.CurrentPage-halfWindow) == 1)
				end := pagination.CurrentPage + halfWindow - bools.ToInt((pagination.CurrentPage+halfWindow) == pagination.NumberOfPages)

				for i := start; i <= end; i++ {
					h.PageSelectorButton("", i, i == pagination.CurrentPage)
				}

				h.PageSelectorEllipsis()
				h.PageSelectorButton("", pagination.NumberOfPages, pagination.NumberOfPages == pagination.CurrentPage)
			case pagination.CurrentPage > pagination.NumberOfPages-pagination.WindowSize:
				h.PageSelectorButton("", 1, pagination.CurrentPage == 1)
				h.PageSelectorEllipsis()
				for i := pagination.NumberOfPages - pagination.WindowSize + 1; i <= pagination.NumberOfPages; i++ {
					h.PageSelectorButton("", i, i == pagination.CurrentPage)
				}
			}
		}

		next := pagination.CurrentPage + bools.ToInt(pagination.CurrentPage < pagination.NumberOfPages)
		h.PageSelectorButton("&raquo;", next, false)

		h.DivEnd()
	}
}
