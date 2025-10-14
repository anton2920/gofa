package http

import (
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/cpu"
	"github.com/anton2920/gofa/debug"
	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/mime/multipart"
	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/session"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace"
)

type Router func(*Response, *Request, session.Session) error

func RequestHandler(w *Response, r *Request, session session.Session, router Router) (err error) {
	t := trace.Begin("")

	defer func() {
		if p := recover(); p != nil {
			r.Error = errors.NewPanic(p)
			err = router(w, r, session)
			trace.End(t)
		}
	}()

	if r.Error == nil {
		switch r.Method {
		case MethodGet:
			if len(r.URL.RawQuery) > 0 {
				err = r.URL.ParseQuery()
			}
		case MethodPost:
			if len(r.Body) > 0 {
				contentType := r.Headers.Get("Content-Type")
				switch {
				case contentType == "application/x-www-form-urlencoded":
					err = url.ParseQuery(&r.Form, bytes.AsString(r.Body))
				case strings.StartsWith(contentType, "multipart/form-data; boundary="):
					err = multipart.ParseFormData(contentType, &r.Form, &r.Files, r.Body)
				}
			}
		}
		if err != nil {
			r.Error = ClientError(err)
		}
	}

	err = router(w, r, session)
	trace.End(t)
	return
}

func RequestsHandler(ws []Response, rs []Request, router Router) {
	t := trace.Begin("")

	const cookie = "Token"

	for i := 0; i < ints.Min(len(ws), len(rs)); i++ {
		w := &ws[i]
		r := &rs[i]

		if r.URL.Path == "/plaintext" {
			const response = "Hello, world!\n"
			switch r.Method {
			default:
				w.WriteString(response)
			case MethodHead:
				w.Headers.Set("Content-Length", "14")
			}
			continue
		}

		start := cpu.ReadPerformanceCounter()
		w.Headers.Set("Content-Type", `text/html; charset="UTF-8"`)
		level := log.LevelDebug

		s := session.Get(r.Cookie(cookie))
		if len(s.Token) == 0 {
			s = session.New(0)
			if debug.Debug {
				w.SetCookieUnsafe(cookie, s.Token, s.Expiry)
			} else {
				w.SetCookie(cookie, s.Token, s.Expiry)
			}
		}

		err := RequestHandler(w, r, s, router)
		if err != nil {
			if (w.StatusCode >= StatusBadRequest) && (w.StatusCode < StatusInternalServerError) {
				level = log.LevelWarn
			} else {
				level = log.LevelError
			}
		}

		if r.Method == MethodHead {
			buffer := w.Arena.NewSlice(ints.Bufsize)
			n := slices.PutInt(buffer, len(w.Body))
			w.Headers.Set("Content-Length", bytes.AsString(buffer[:n]))
			w.Body = w.Body[:0]
		}

		if r.Headers.Get("Connection") == "close" {
			w.Headers.Set("Connection", "close")
		}

		end := cpu.ReadPerformanceCounter()
		elapsed := end - start

		log.Logf(level, "[%21s] %7s %s -> %v (%v), %4dus", strings.And(r.RemoteAddr, r.Headers.Get("X-Forwarded-For")), r.Method, r.URL.Path, w.StatusCode, err, elapsed.ToMicroseconds())
	}

	trace.End(t)
}

func ConnectionHandler(c *Conn, router Router) {
	const pipeline = 64

	rs := make([]Request, pipeline)
	ws := make([]Response, pipeline)

	for !c.Closed {
		n, err := c.ReadRequests(rs)
		if err != nil {
			if n == 0 {
				break
			}
			log.Errorf("Failed to read HTTP requests: %v", err)
			break
		}

		if n > 0 {
			RequestsHandler(ws[:n], rs[:n], router)

			n, err = c.WriteResponses(ws[:n])
			if err != nil {
				log.Errorf("Failed to write HTTP responses: %v", err)
				break
			}
		}
	}

	if err := c.Close(); err != nil {
		log.Warnf("Failed to close HTTP connection: %v", err)
	}

	c.ConnPool.Put(c)
}
