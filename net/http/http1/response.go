package http1

import (
	"unsafe"

	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/prof"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/time"
)

var StatusLines = [...]string{
	0:                                "HTTP/1.1 200 OK\r\n",
	http.StatusOK:                    "HTTP/1.1 200 OK\r\n",
	http.StatusSeeOther:              "HTTP/1.1 303 See Other\r\n",
	http.StatusBadRequest:            "HTTP/1.1 400 Bad Request\r\n",
	http.StatusUnauthorized:          "HTTP/1.1 401 Unauthorized\r\n",
	http.StatusForbidden:             "HTTP/1.1 403 Forbidden\r\n",
	http.StatusNotFound:              "HTTP/1.1 404 Not Found\r\n",
	http.StatusMethodNotAllowed:      "HTTP/1.1 405 Method Not Allowed\r\n",
	http.StatusRequestTimeout:        "HTTP/1.1 408 Request Timeout\r\n",
	http.StatusConflict:              "HTTP/1.1 409 Conflict\r\n",
	http.StatusRequestEntityTooLarge: "HTTP/1.1 413 Request Entity Too Large\r\n",
	http.StatusInternalServerError:   "HTTP/1.1 500 Internal Server Error\r\n",
}

func FillError(ctx *http.Context, err error, dateBuf []byte) {
	p := prof.Begin("")

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

	w.WriteString(http.Status2Reason[w.StatusCode])
	w.WriteString(`: `)
	w.WriteString(message)
	w.WriteString("\r\n")

	w.Headers.Set("Connection", "close")
	FillResponses(ctx, unsafe.Slice(&w, 1), dateBuf)

	prof.End(p)
}

func FillResponses(ctx *http.Context, ws []http.Response, dateBuf []byte) {
	p := prof.Begin("")

	for i := 0; i < len(ws); i++ {
		w := &ws[i]

		ctx.ResponseBuffer = append(ctx.ResponseBuffer, StatusLines[w.StatusCode]...)

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
			ctx.ResponseBuffer = append(ctx.ResponseBuffer, "Content-Type: text/plain; charset=\"UTF-8\"\r\n"...)
		}

		if !w.Headers.Has("Content-Length") {
			lengthBuf := make([]byte, 20)
			n := slices.PutInt(lengthBuf, len(w.Body))

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
		ctx.ResponseBuffer = append(ctx.ResponseBuffer, w.Body...)
		w.Reset()
	}

	prof.End(p)
}
