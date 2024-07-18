package http1

import (
	"unsafe"

	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/time"
)

func FillError(ctx *http.Context, err error, dateBuf []byte) {
	var w http.Response
	var message string

	switch err := err.(type) {
	default:
		w.StatusCode = http.StatusInternalServerError
		message = err.Error()
	case http.Error:
		w.StatusCode = err.StatusCode
		message = err.DisplayMessage
	}

	w.AppendString(http.Status2Reason[w.StatusCode])
	w.AppendString(`: `)
	w.WriteString(message)
	w.AppendString("\r\n")

	w.Headers.Set("Connection", "close")
	FillResponses(ctx, unsafe.Slice(&w, 1), dateBuf)
}

func FillResponses(ctx *http.Context, ws []http.Response, dateBuf []byte) {
	for i := 0; i < len(ws); i++ {
		w := &ws[i]

		ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("HTTP/1.1"), syscall.Iovec(" "), syscall.Iovec(http.Status2String[w.StatusCode]), syscall.Iovec(" "), syscall.Iovec(http.Status2Reason[w.StatusCode]), syscall.Iovec("\r\n"))

		if !w.Headers.Has("Date") {
			if dateBuf == nil {
				dateBuf = make([]byte, 31)
				time.PutTmRFC822(dateBuf, time.ToTm(time.Unix()))
			}

			ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("Date: "), syscall.IovecForByteSlice(dateBuf), syscall.Iovec("\r\n"))
		}

		if !w.Headers.Has("Server") {
			ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("Server: gofa/http\r\n"))
		}

		if !w.Headers.Has("Content-Type") {
			if http.ContentTypeHTML(w.Bodies) {
				ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("Content-Type: text/html; charset=\"UTF-8\"\r\n"))
			} else {
				ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("Content-Type: text/plain; charset=\"UTF-8\"\r\n"))
			}
		}

		if !w.Headers.Has("Content-Length") {
			var length int
			for i := 0; i < len(w.Bodies); i++ {
				length += int(len(w.Bodies[i]))
			}

			lengthBuf := w.Arena.NewSlice(20)
			n := slices.PutInt(lengthBuf, length)

			ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("Content-Length: "), syscall.IovecForByteSlice(lengthBuf[:n]), syscall.Iovec("\r\n"))
		}

		for i := 0; i < len(w.Headers.Keys); i++ {
			key := w.Headers.Keys[i]
			ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec(key), syscall.Iovec(": "))
			for j := 0; j < len(w.Headers.Values[i]); j++ {
				value := w.Headers.Values[i][j]
				if j > 0 {
					ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec(", "))
				}
				ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec(value))
			}
			ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("\r\n"))
		}

		ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("\r\n"))
		ctx.ResponseIovs = append(ctx.ResponseIovs, w.Bodies...)
		w.Reset()
	}
}
