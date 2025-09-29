package http

import (
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/cpu"
	"github.com/anton2920/gofa/debug"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/mime/multipart"
	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/session"
	"github.com/anton2920/gofa/strings"
)

type Router func(*Response, *Request, session.Session) error

func RequestHandler(w *Response, r *Request, session session.Session, router Router) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = error.NewPanic(p)
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
			CilentError(err)
		}
	}

	return router(w, r, session)
}

func RequestsHandler(ws []Response, rs []Request, router Router) {
	const cookie = "Token"

	for i := 0; i < ints.Min(len(ws), len(rs)); i++ {
		w := &ws[i]
		r := &rs[i]

		start := cpu.GetPerformanceCounter()
		w.Headers.Set("Content-Type", `text/html; charset="UTF-8"`)
		level := log.LevelDebug

		session := session.Get(r.Cookie(cookie))
		if len(session.Token) == 0 {
			session = session.New(0)
			if debug.Debug {
				w.SetCookieUnsafe(cookie, session.Token, int(session.Expiry))
			} else {
				w.SetCookie(cookie, session.Token, int(session.Expiry))
			}
		}

		err := router(w, r, session)
		if err != nil {
			if (w.StatusCode >= StatusBadRequest) && (w.StatusCode < StatusInternalServerError) {
				level = log.LevelWarn
			} else {
				level = log.LevelError
			}
		}

		if r.Headers.Get("Connection") == "close" {
			w.Headers.Set("Connection", "close")
		}

		end := cpu.GetPerformanceCounter()
		elapsed := end - start

		log.Logf(level, "[%21s] %7s %s -> %v (%v), %4dus", strings.And(r.Address, r.Headers.Get("X-Forwarded-For")), r.Method, r.URL.Path, w.StatusCode, err, elapsed.ToUsec())
	}
}

func ConnectionHandler(c *Conn, handler func([]Response, []Request)) {
	const pipeline = 64

	rs := make([]Request, pipeline)
	ws := make([]Response, pipeline)

	for !c.Closed {
		n, err := ReadRequests(c, rs)
		if err != nil {
			log.Errorf("Failed to read HTTP requests: %v", err)
			break
		} else if n == 0 {
			break
		}

		handler(ws[:n], rs[:n])

		n, err = WriteResponses(c, ws[:n])
		if err != nil {
			log.Errorf("Failed to write HTTP responses: %v", err)
			break
		} else if (n > 0) && (n < len(ws)) {
			log.Errorf("Failed to write all HTTP responses: %v", err)
			break
		}
	}

	if err := c.Close(); err != nil {
		log.Warnf("Failed to close HTTP connection: %v", err)
	}
}
