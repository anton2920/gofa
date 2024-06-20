package http

import "github.com/anton2920/gofa/syscall"

type Headers struct {
	Values []syscall.Iovec

	OmitDate          bool
	OmitServer        bool
	OmitContentType   bool
	OmitContentLength bool
}
