package http1

import (
	"unsafe"

	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/slices"
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

		/* TODO(anton2920): prepare an array of status lines. */
		ctx.ResponseBuffer = append(ctx.ResponseBuffer, "HTTP/1.1"...)
		ctx.ResponseBuffer = append(ctx.ResponseBuffer, " "...)
		ctx.ResponseBuffer = append(ctx.ResponseBuffer, http.Status2String[w.StatusCode]...)
		ctx.ResponseBuffer = append(ctx.ResponseBuffer, " "...)
		ctx.ResponseBuffer = append(ctx.ResponseBuffer, http.Status2Reason[w.StatusCode]...)
		ctx.ResponseBuffer = append(ctx.ResponseBuffer, "\r\n"...)

		if !w.Headers.Has("Date") {
			if dateBuf == nil {
				dateBuf = make([]byte, time.RFC822Len)
				time.PutTmRFC822(dateBuf, time.ToTm(time.Unix()))
			}
			ctx.ResponseBuffer = append(ctx.ResponseBuffer, "Date: "...)
			ctx.ResponseBuffer = append(ctx.ResponseBuffer, dateBuf...)
			ctx.ResponseBuffer = append(ctx.ResponseBuffer, "\r\n"...)
		}

		if !w.Headers.Has("Server") {
			ctx.ResponseBuffer = append(ctx.ResponseBuffer, "Server: gofa/http\r\n"...)
		}

		if !w.Headers.Has("Content-Type") {
			var contentType string
			if http.ContentTypeHTML(w.Bodies) {
				contentType = `text/html; charset="UTF-8"`
			} else {
				contentType = `text/plain; charset="UTF-8"`
			}
			ctx.ResponseBuffer = append(ctx.ResponseBuffer, "Content-Type: "...)
			ctx.ResponseBuffer = append(ctx.ResponseBuffer, contentType...)
			ctx.ResponseBuffer = append(ctx.ResponseBuffer, "\r\n"...)
		}

		if !w.Headers.Has("Content-Length") {
			var length int
			for i := 0; i < len(w.Bodies); i++ {
				length += int(len(w.Bodies[i]))
			}

			lengthBuf := w.Arena.NewSlice(20)
			n := slices.PutInt(lengthBuf, length)

			ctx.ResponseBuffer = append(ctx.ResponseBuffer, "Content-Length: "...)
			ctx.ResponseBuffer = append(ctx.ResponseBuffer, lengthBuf[:n]...)
			ctx.ResponseBuffer = append(ctx.ResponseBuffer, "\r\n"...)
		}

		for i := 0; i < len(w.Headers.Keys); i++ {
			key := w.Headers.Keys[i]
			ctx.ResponseBuffer = append(ctx.ResponseBuffer, key...)
			ctx.ResponseBuffer = append(ctx.ResponseBuffer, ": "...)
			for j := 0; j < len(w.Headers.Values[i]); j++ {
				value := w.Headers.Values[i][j]
				if j > 0 {
					ctx.ResponseBuffer = append(ctx.ResponseBuffer, ","...)
				}
				ctx.ResponseBuffer = append(ctx.ResponseBuffer, value...)
			}
			ctx.ResponseBuffer = append(ctx.ResponseBuffer, "\r\n"...)
		}

		ctx.ResponseBuffer = append(ctx.ResponseBuffer, "\r\n"...)
		for i := 0; i < len(w.Bodies); i++ {
			ctx.ResponseBuffer = append(ctx.ResponseBuffer, w.Bodies[i]...)
		}
		w.Reset()
	}
}
