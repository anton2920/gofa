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
	"github.com/anton2920/gofa/strings"
)

type Router func(*Response, *Request, session.Session) error

func RequestHandler(w *Response, r *Request, session session.Session, router Router) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = errors.NewPanic(p)
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
			ClientError(err)
		}
	}

	return router(w, r, session)
}

func RequestsHandler(ws []Response, rs []Request, router Router) {
	const cookie = "Token"

	for i := 0; i < ints.Min(len(ws), len(rs)); i++ {
		w := &ws[i]
		r := &rs[i]

		start := cpu.ReadPerformanceCounter()
		w.Headers.Set("Content-Type", `text/html; charset="UTF-8"`)
		level := log.LevelDebug

		session := session.Get(r.Cookie(cookie))
		if len(session.Token) == 0 {
			session = session.New(0)
			if debug.Debug {
				w.SetCookieUnsafe(cookie, session.Token, session.Expiry)
			} else {
				w.SetCookie(cookie, session.Token, session.Expiry)
			}
		}

		err := RequestHandler(w, r, session, router)
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

		end := cpu.ReadPerformanceCounter()
		elapsed := end - start

		log.Logf(level, "[%21s] %7s %s -> %v (%v), %4dus", strings.And(r.RemoteAddr, r.Headers.Get("X-Forwarded-For")), r.Method, r.URL.Path, w.StatusCode, err, elapsed.ToMicroseconds())
	}
}

func ConnectionHandler(l *Listener, c *Conn, handler func([]Response, []Request, Router), router Router) {
	const pipeline = 64

	rs := make([]Request, pipeline)
	ws := make([]Response, pipeline)

	for !c.Closed {
		n, err := c.ReadRequests(rs)
		if err != nil {
			log.Errorf("Failed to read HTTP requests: %v", err)
			break
		} else if n == 0 {
			break
		}

		handler(ws[:n], rs[:n], router)

		n, err = c.WriteResponses(ws[:n])
		if err != nil {
			log.Errorf("Failed to write HTTP responses: %v", err)
			break
		}
	}

	if err := c.Close(); err != nil {
		log.Warnf("Failed to close HTTP connection: %v", err)
	}

	l.ConnPool.Put(c)
}
