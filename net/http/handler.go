package http

import (
	"unsafe"

	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/cpu"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/mime/multipart"
	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/pointers"
	"github.com/anton2920/gofa/session"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace/trace_"
)

type Router func(*context.Context, *Response, *Request) bool

const Pipeline = 16

func RequestHandler(ctx *context.Context, w *Response, r *Request, router Router) bool {
	t := trace_.Begin("")

	if ctx.OK() {
		switch r.Method {
		case MethodGet:
			if len(r.URL.RawQuery) > 0 {
				r.URL.ParseQuery(ctx)
			}
		case MethodPost:
			if len(r.Body) > 0 {
				contentType := r.Headers.Get("Content-Type")
				switch {
				case contentType == "application/x-www-form-urlencoded":
					url.ParseQuery(ctx, &r.Form, bytes.AsString(r.Body))
				case strings.StartsWith(contentType, "multipart/form-data; boundary="):
					multipart.ParseFormData(ctx, contentType, &r.Form, &r.Files, r.Body)
				}
			}
		}
	}

	ok := router(ctx, (*Response)(pointers.UnsafeNoescape(unsafe.Pointer(w))), (*Request)(pointers.UnsafeNoescape(unsafe.Pointer(r))))
	trace_.End(t)
	return ok
}

func RequestsHandler(ctx *context.Context, ws []Response, rs []Request, router Router) {
	t := trace_.Begin("")

	const cookie = "Token"

	for i := 0; i < len(rs); i++ {
		w := &ws[i]
		r := &rs[i]

		if r.URL.Path == "/plaintext" {
			const response = "Hello, world!\n"
			switch r.Method {
			default:
				w.WriteString(response)
				//case MethodHead:
				//	w.Headers.Set("Content-Length", "14")
			}
			continue
		}

		start := cpu.ReadPerformanceCounter()
		w.Headers.Set("Content-Type", `text/html; charset="UTF-8"`)
		level := log.LevelInfo

		/* TODO(anton2920): store session.Customization on client. */
		r.Session = session.Get(r.Cookie(cookie))
		/*
			if len(r.Token) == 0 {
				r.Session = session.New(0)
				if debug.Debug {
					w.SetCookieUnsafe(cookie, r.Token, r.Expiry)
				} else {
					w.SetCookie(cookie, r.Token, r.Expiry)
				}
			}
		*/

		if !RequestHandler(ctx, w, r, router) {
			if (w.Status >= StatusBadRequest) && (w.Status < StatusInternalServerError) {
				level = log.LevelWarn
			} else {
				level = log.LevelError
			}
		}

		if r.Method == MethodHead {
			buffer := ctx.Arena.PushByteArray(ints.Bufsize)
			n := slices.PutInt(buffer, len(w.Body))
			w.Headers.Set("Content-Length", bytes.AsString(buffer[:n]))
			w.Body = w.Body[:0]
		}

		if r.Headers.Get("Connection") == "close" {
			w.Headers.Set("Connection", "close")
		}

		end := cpu.ReadPerformanceCounter()
		elapsed := end - start

		_, _ = level, elapsed
		// log_.Println(ctx.Log.Log(level, time_.NowInNanoseconds()).S("[").W(21).S(strings.Or(r.Headers.Get("X-Forwarded-For"), r.RemoteAddr)).S("] ").W(7).S(r.Method).S(" ").S(string(r.URL.Path)).S(" -> ").S(w.Status.String()).S(" (").S(ctx.Error()).S("), ").W(4).D64(elapsed.ToMicrosecondsTruncated()).S("us"))
	}

	trace_.End(t)
}
