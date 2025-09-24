package html

import "github.com/anton2920/gofa/bools"

type Pagination struct {
	CurrentPage   int
	ItemsPerPage  int
	NumberOfPages int
	WindowSize    int
}

func GetNumberOfPages(perPage int, ntotal int) int {
	return (ntotal / perPage) + bools.ToInt(ntotal%perPage > 0)
}

func (h *HTML) PaginationButton(page int, active bool) {
	attrs := h.Theme.PaginationButton
	if active {
		attrs = h.MergeAttributes(attrs, h.Theme.PaginationButtonActive)
	}
	h.WithoutTheme().Button(h.Itoa(page), attrs, Attributes{Name: "Page", FormNoValidate: true})
}

func (h *HTML) PaginationEllipsis() {
	h.WithoutTheme().Button("...", h.Theme.PaginationButton, Attributes{Disabled: true})
}

func (h *HTML) Pagination(pagination *Pagination, attrs ...Attributes) {
	if pagination.NumberOfPages > 1 {
		h.DivBegin(attrs...)

		h.String(` <button class="join-item btn" name="Page" value="`)
		if pagination.CurrentPage == 0 {
			h.Int(pagination.CurrentPage)
		} else {
			h.Int(pagination.CurrentPage - 1)
		}
		h.String(`">«</button>`)

		if pagination.WindowSize > pagination.NumberOfPages {
			pagination.WindowSize = pagination.NumberOfPages - 1
		}

		if pagination.NumberOfPages == 2 {
			for i := 0; i < pagination.NumberOfPages; i++ {
				h.PaginationButton(i, i == pagination.CurrentPage)
			}
		} else {
			switch {
			case pagination.CurrentPage < pagination.WindowSize-1:
				for i := 0; i < pagination.WindowSize; i++ {
					h.PaginationButton(i, i == pagination.CurrentPage)
				}
				h.PaginationEllipsis()
				h.PaginationButton(pagination.NumberOfPages-1, pagination.NumberOfPages-1 == pagination.CurrentPage)
			default:
				h.PaginationButton(0, 0 == pagination.CurrentPage)
				h.PaginationEllipsis()

				start := pagination.CurrentPage - pagination.WindowSize/2 + bools.ToInt((pagination.CurrentPage-pagination.WindowSize/2) == 0)
				end := pagination.CurrentPage + pagination.WindowSize/2 - bools.ToInt((pagination.CurrentPage+pagination.WindowSize/2) == pagination.NumberOfPages-1)
				for i := start; i <= end; i++ {
					h.PaginationButton(i, i == pagination.CurrentPage)
				}

				h.PaginationEllipsis()
				h.PaginationButton(pagination.NumberOfPages-1, pagination.NumberOfPages-1 == pagination.CurrentPage)
			case pagination.CurrentPage > pagination.NumberOfPages-pagination.WindowSize:
				h.PaginationButton(0, 0 == pagination.CurrentPage)
				h.PaginationEllipsis()
				for i := pagination.NumberOfPages - pagination.WindowSize; i < pagination.NumberOfPages; i++ {
					h.PaginationButton(i, i == pagination.CurrentPage)
				}
			}
		}

		h.String(` <button class="join-item btn" name="Page" value="`)
		if pagination.CurrentPage == pagination.NumberOfPages-1 {
			h.Int(pagination.CurrentPage)
		} else {
			h.Int(pagination.CurrentPage + 1)
		}
		h.String(`">»</button>`)

		h.DivEnd()
	}
}
