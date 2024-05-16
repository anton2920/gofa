package http

import (
	"fmt"
	"log"
	"strconv"
	"unsafe"

	"github.com/anton2920/gofa/arena"
	"github.com/anton2920/gofa/buffer"
	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/event"
	"github.com/anton2920/gofa/net/tcp"
	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/time"
)

type State int

const (
	StateRequestLine State = iota
	StateHeader
	StateBody

	StateUnknown
	StateDone
)

type Request struct {
	Arena arena.Arena

	RemoteAddr string

	Method string
	URL    url.URL
	Proto  string

	Headers []string
	Body    []byte

	Form url.Values
}

type RequestParser struct {
	State State
	Pos   int

	ContentLength int
}

type Status int

const (
	StatusOK                    Status = 200
	StatusSeeOther                     = 303
	StatusBadRequest                   = 400
	StatusUnauthorized                 = 401
	StatusForbidden                    = 403
	StatusNotFound                     = 404
	StatusMethodNotAllowed             = 405
	StatusRequestTimeout               = 408
	StatusConflict                     = 409
	StatusRequestEntityTooLarge        = 413
	StatusInternalServerError          = 500
)

var Status2String = [...]string{
	0:                           "200",
	StatusOK:                    "200",
	StatusSeeOther:              "303",
	StatusBadRequest:            "400",
	StatusUnauthorized:          "401",
	StatusForbidden:             "403",
	StatusNotFound:              "404",
	StatusMethodNotAllowed:      "405",
	StatusRequestTimeout:        "408",
	StatusConflict:              "409",
	StatusRequestEntityTooLarge: "413",
	StatusInternalServerError:   "500",
}

var Status2Reason = [...]string{
	0:                           "OK",
	StatusOK:                    "OK",
	StatusSeeOther:              "See Other",
	StatusBadRequest:            "Bad Request",
	StatusUnauthorized:          "Unauthorized",
	StatusForbidden:             "Forbidden",
	StatusNotFound:              "Not Found",
	StatusMethodNotAllowed:      "Method Not Allowed",
	StatusRequestTimeout:        "Request Timeout",
	StatusConflict:              "Conflict",
	StatusRequestEntityTooLarge: "Request Entity Too Large",
	StatusInternalServerError:   "Internal Server Error",
}

type Headers struct {
	Values []syscall.Iovec

	OmitDate          bool
	OmitServer        bool
	OmitContentType   bool
	OmitContentLength bool
}

type Response struct {
	Arena arena.Arena

	StatusCode Status
	Headers    Headers
	Bodies     []syscall.Iovec
}

type Error struct {
	StatusCode     Status
	DisplayMessage string
	LogError       error
}

type Context struct {
	/* Check must be the same as the last pointer's bit, if context is in use. */
	Check int32

	Connection    int32
	ClientAddress string

	RequestPendingBytes int
	RequestParser       RequestParser
	RequestBuffer       buffer.Circular

	ResponseIovs []syscall.Iovec
	ResponsePos  int

	/* DateRFC822 could be set by client to reduce unnecessary syscalls and date formatting. */
	DateRFC822 []byte

	/* Optional event queue this client is attached to. */
	EventQueue *event.Queue
}

var (
	UnauthorizedError = Error{StatusCode: StatusUnauthorized, DisplayMessage: "whoops... You have to sign in to see this page", LogError: errors.New("whoops... You have to sign in to see this page")}
	ForbiddenError    = Error{StatusCode: StatusForbidden, DisplayMessage: "whoops... Your permissions are insufficient", LogError: errors.New("whoops... Your permissions are insufficient")}
)

const HTMLHeader = `<!DOCTYPE html>`

var (
	HTMLQuot = "&#34;" // shorter than "&quot;"
	HTMLApos = "&#39;" // shorter than "&apos;" and apos was not in HTML until HTML5
	HTMLAmp  = "&amp;"
	HTMLLt   = "&lt;"
	HTMLGt   = "&gt;"
	HTMLNull = "\uFFFD"
)

func (s Status) String() string {
	return Status2String[s]
}

func (rp *RequestParser) Parse(request string, r *Request) (int, error) {
	var err error

	rp.State = StateRequestLine
	rp.ContentLength = 0
	rp.Pos = 0

	for rp.State != StateDone {
		switch rp.State {
		default:
			log.Panicf("Unknown  parser state %d", rp.State)
		case StateUnknown:
			if len(request[rp.Pos:]) < 2 {
				return 0, nil
			}
			if request[rp.Pos:rp.Pos+2] == "\r\n" {
				rp.Pos += len("\r\n")

				if rp.ContentLength != 0 {
					rp.State = StateBody
				} else {
					rp.State = StateDone
				}
			} else {
				rp.State = StateHeader
			}
		case StateRequestLine:
			lineEnd := strings.FindChar(request[rp.Pos:], '\r')
			if lineEnd == -1 {
				return 0, nil
			}

			sp := strings.FindChar(request[rp.Pos:rp.Pos+lineEnd], ' ')
			if sp == -1 {
				return 0, fmt.Errorf("expected method, found %q", request[rp.Pos:])
			}
			r.Method = request[rp.Pos : rp.Pos+sp]
			rp.Pos += len(r.Method) + 1
			lineEnd -= len(r.Method) + 1

			uriEnd := strings.FindChar(request[rp.Pos:rp.Pos+lineEnd], ' ')
			if uriEnd == -1 {
				return 0, fmt.Errorf("expected space after URI, found %q", request[rp.Pos:lineEnd])
			}

			queryStart := strings.FindChar(request[rp.Pos:rp.Pos+uriEnd], '?')
			if queryStart != -1 {
				r.URL.Path = request[rp.Pos : rp.Pos+queryStart]
				r.URL.Query = request[rp.Pos+queryStart+1 : rp.Pos+uriEnd]
			} else {
				r.URL.Path = request[rp.Pos : rp.Pos+uriEnd]
				r.URL.Query = ""
			}
			rp.Pos += len(r.URL.Path) + len(r.URL.Query) + 1
			lineEnd -= len(r.URL.Path) + len(r.URL.Query) + 1

			if request[rp.Pos:rp.Pos+len("HTTP/")] != "HTTP/" {
				return 0, fmt.Errorf("expected version prefix, found %q", request[rp.Pos:rp.Pos+lineEnd])
			}
			r.Proto = request[rp.Pos : rp.Pos+lineEnd]

			rp.Pos += len(r.Proto) + len("\r\n")
			rp.State = StateUnknown
		case StateHeader:
			lineEnd := strings.FindChar(request[rp.Pos:], '\r')
			if lineEnd == -1 {
				return 0, nil
			}
			header := request[rp.Pos : rp.Pos+lineEnd]
			r.Headers = append(r.Headers, header)

			if strings.StartsWith(header, "Content-Length: ") {
				header = header[len("Content-Length: "):]
				rp.ContentLength, err = strconv.Atoi(header)
				if err != nil {
					return 0, fmt.Errorf("failed to parse Content-Length value: %w", err)
				}
			}

			rp.Pos += len(header) + len("\r\n")
			rp.State = StateUnknown
		case StateBody:
			if len(request[rp.Pos:]) < rp.ContentLength {
				return 0, nil
			}

			r.Body = unsafe.Slice(unsafe.StringData(request[rp.Pos:]), rp.ContentLength)
			rp.Pos += len(r.Body)
			rp.State = StateDone
		}
	}

	return rp.Pos, nil
}

func (r *Request) Cookie(name string) string {
	for i := 0; i < len(r.Headers); i++ {
		header := r.Headers[i]
		if strings.StartsWith(header, "Cookie: ") {
			cookie := header[len("Cookie: "):]
			if strings.StartsWith(cookie, name) {
				cookie = cookie[len(name):]
				if cookie[0] != '=' {
					return ""
				}
				return cookie[1:]
			}

		}
	}

	return ""
}

func (r *Request) ParseForm() error {
	var err error

	if len(r.Form) != 0 {
		return nil
	}

	query := unsafe.String(unsafe.SliceData(r.Body), len(r.Body))
	for query != "" {
		var key string
		key, query, _ = strings.Cut(query, "&")
		if strings.FindChar(key, ';') != -1 {
			err = errors.New("invalid semicolon separator in query")
			continue
		}
		if key == "" {
			continue
		}
		key, value, _ := strings.Cut(key, "=")

		keyBuffer := r.Arena.NewSlice(len(key))
		n, ok := url.QueryDecode(keyBuffer, key)
		if !ok {
			if err == nil {
				err = errors.New("invalid key")
			}
			continue
		}
		key = unsafe.String(unsafe.SliceData(keyBuffer), n)

		valueBuffer := r.Arena.NewSlice(len(value))
		n, ok = url.QueryDecode(valueBuffer, value)
		if !ok {
			if err == nil {
				err = errors.New("invalid value")
			}
			continue
		}
		value = unsafe.String(unsafe.SliceData(valueBuffer), n)

		r.Form.Add(key, value)
	}

	return err
}

func (w *Response) Append(b []byte) {
	w.Bodies = append(w.Bodies, syscall.IovecForByteSlice(b))
}

func (w *Response) AppendString(s string) {
	w.Bodies = append(w.Bodies, syscall.Iovec(s))
}

func (w *Response) DelCookie(name string) {
	const finisher = "=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Strict"

	cookie := w.Arena.NewSlice(len(name) + len(finisher))

	var n int
	n += copy(cookie[n:], name)
	n += copy(cookie[n:], finisher)

	w.SetHeaderUnsafe("Set-Cookie", unsafe.String(unsafe.SliceData(cookie), n))
}

func (w *Response) SetCookie(name, value string, expiry int64) {
	const secure = "; HttpOnly; Secure; SameSite=Strict"
	const expires = "; Expires="
	const path = "; Path=/"
	const eq = "="

	cookie := w.Arena.NewSlice(len(name) + len(eq) + len(value) + len(path) + len(expires) + time.RFC822Len + len(secure))

	var n int
	n += copy(cookie[n:], name)
	n += copy(cookie[n:], eq)
	n += copy(cookie[n:], value)
	n += copy(cookie[n:], path)
	n += copy(cookie[n:], expires)
	n += time.PutTmRFC822(cookie[n:], time.ToTm(int(expiry)))
	n += copy(cookie[n:], secure)

	w.SetHeaderUnsafe("Set-Cookie", unsafe.String(unsafe.SliceData(cookie), n))
}

/* SetCookieUnsafe is useful for debugging purposes. It's also more compatible with older browsers. */
func (w *Response) SetCookieUnsafe(name, value string, expiry int64) {
	const expires = "; Expires="
	const path = "; Path=/"
	const eq = "="

	cookie := w.Arena.NewSlice(len(name) + len(eq) + len(value) + len(path) + len(expires) + time.RFC822Len)

	var n int
	n += copy(cookie[n:], name)
	n += copy(cookie[n:], eq)
	n += copy(cookie[n:], value)
	n += copy(cookie[n:], path)
	n += copy(cookie[n:], expires)
	n += time.PutTmRFC822(cookie[n:], time.ToTm(int(expiry)))

	w.SetHeaderUnsafe("Set-Cookie", unsafe.String(unsafe.SliceData(cookie), n))
}

/* SetHeaderUnsafe sets new 'value' for 'header' relying on that memory lives long enough. */
func (w *Response) SetHeaderUnsafe(header string, value string) {
	switch header {
	case "Date":
		w.Headers.OmitDate = true
	case "Server":
		w.Headers.OmitServer = true
	case "ContentType":
		w.Headers.OmitContentType = true
	case "ContentLength":
		w.Headers.OmitContentLength = true
	}

	for i := 0; i < len(w.Headers.Values); i += 4 {
		key := w.Headers.Values[i]
		if header == string(key) {
			w.Headers.Values[i+2] = syscall.Iovec(value)
			return
		}
	}

	w.Headers.Values = append(w.Headers.Values, syscall.Iovec(header), syscall.Iovec(": "), syscall.Iovec(value), syscall.Iovec("\r\n"))
}

func (w *Response) Redirect(path string, code Status) {
	pathBuf := w.Arena.NewSlice(len(path))
	copy(pathBuf, path)

	w.SetHeaderUnsafe("Location", unsafe.String(unsafe.SliceData(pathBuf), len(pathBuf)))
	w.Bodies = w.Bodies[:0]
	w.StatusCode = code
}

func (w *Response) RedirectID(prefix string, id int, code Status) {
	buffer := w.Arena.NewSlice(len(prefix) + 20)
	n := copy(buffer, prefix)
	n += slices.PutInt(buffer[n:], id)

	w.SetHeaderUnsafe("Location", unsafe.String(unsafe.SliceData(buffer), n))
	w.Bodies = w.Bodies[:0]
	w.StatusCode = code
}

func (w *Response) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	buffer := w.Arena.NewSlice(len(b))
	copy(buffer, b)
	w.Append(buffer)

	return len(b), nil
}

/* WriteHTML writes to w the escaped HTML equivalent of the plain text data b. */
func (w *Response) WriteHTML(b []byte) {
	last := 0
	for i, c := range b {
		var html string
		switch c {
		case '\000':
			html = HTMLNull
		case '"':
			html = HTMLQuot
		case '\'':
			html = HTMLApos
		case '&':
			html = HTMLAmp
		case '<':
			html = HTMLLt
		case '>':
			html = HTMLGt
		default:
			continue
		}
		w.Write(b[last:i])
		w.AppendString(html)
		last = i + 1
	}
	w.Write(b[last:])
}

func (w *Response) WriteInt(i int) (int, error) {
	buffer := w.Arena.NewSlice(20)
	n := slices.PutInt(buffer, i)
	w.Append(buffer[:n])
	return n, nil
}

func (w *Response) WriteString(s string) (int, error) {
	return w.Write(unsafe.Slice(unsafe.StringData(s), len(s)))
}

func (w *Response) WriteHTMLString(s string) {
	w.WriteHTML(unsafe.Slice(unsafe.StringData(s), len(s)))
}

func BadRequest(format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{StatusCode: StatusBadRequest, DisplayMessage: message, LogError: errors.WrapWithTrace(errors.New(message), 2)}
}

func NotFound(format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{StatusCode: StatusNotFound, DisplayMessage: message, LogError: errors.WrapWithTrace(errors.New(message), 2)}
}

func Conflict(format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return Error{StatusCode: StatusConflict, DisplayMessage: message, LogError: errors.WrapWithTrace(errors.New(message), 2)}
}

func ClientError(err error) Error {
	return Error{StatusCode: StatusBadRequest, DisplayMessage: "whoops... Something went wrong. Please reload this page or try again later", LogError: errors.WrapWithTrace(err, 2)}
}

func ServerError(err error) Error {
	return Error{StatusCode: StatusInternalServerError, DisplayMessage: "whoops... Something went wrong. Please try again later", LogError: errors.WrapWithTrace(err, 2)}
}

func (e Error) Error() string {
	if e.LogError == nil {
		return "<nil>"
	}
	return e.LogError.Error()
}

func NewContext(c int32, addr tcp.SockAddrIn, bufferSize int) (*Context, error) {
	rb, err := buffer.NewCircular(bufferSize)
	if err != nil {
		return nil, err
	}

	ctx := new(Context)
	ctx.Connection = c
	ctx.RequestBuffer = rb
	ctx.ResponseIovs = make([]syscall.Iovec, 0, 1024)

	buffer := make([]byte, 21)
	n := tcp.PutAddress(buffer, addr.Addr, addr.Port)
	ctx.ClientAddress = string(buffer[:n])

	return ctx, nil
}

func ContextFromEvent(event event.Event) (*Context, bool) {
	if event.UserData == nil {
		return nil, false
	}
	uptr := uintptr(event.UserData)

	check := uptr & 0x1
	ctx := (*Context)(unsafe.Pointer(uptr - check))
	ctx.RequestPendingBytes = event.Available

	return ctx, ctx.Check == int32(check)
}

func (ctx *Context) Reset() {
	ctx.Check = 1 - ctx.Check
	ctx.RequestBuffer.Reset()
	ctx.ResponsePos = 0
	ctx.ResponseIovs = ctx.ResponseIovs[:0]
}

func FreeContext(ctx *Context) {
	ctx.Reset()
	buffer.FreeCircular(&ctx.RequestBuffer)
}

func Accept(l int32, bufferSize int) (*Context, error) {
	var addr tcp.SockAddrIn
	var addrLen uint32 = uint32(unsafe.Sizeof(addr))

	c, err := syscall.Accept(l, (*syscall.Sockaddr)(unsafe.Pointer(&addr)), &addrLen)
	if err != nil {
		return nil, fmt.Errorf("failed to accept incoming connection: %w", err)
	}

	ctx, err := NewContext(c, addr, bufferSize)
	if err != nil {
		syscall.Close(c)
		return nil, fmt.Errorf("failed to create new  context: %w", err)
	}

	return ctx, nil
}

func AddClientToQueue(q *event.Queue, ctx *Context, request event.Request, trigger event.Trigger) error {
	/* TODO(anton2920): switch to pinning inside platform methods. */
	q.Pinner.Pin(ctx)
	ctx.EventQueue = q
	return q.AddSocket(ctx.Connection, request, trigger, unsafe.Pointer(uintptr(unsafe.Pointer(ctx))|uintptr(ctx.Check)))
}

/* ReadRequests reads data from socket and parses  requests. Returns the number of requests parsed. */
func ReadRequests(ctx *Context, rs []Request) (int, error) {
	usesQ := ctx.EventQueue != nil
	parser := &ctx.RequestParser
	rBuf := &ctx.RequestBuffer

	if (!usesQ) || (ctx.RequestPendingBytes > 0) {
		if rBuf.RemainingSpace() == 0 {
			return 0, errors.New("no space left in the buffer")
		}
		n, err := syscall.Read(ctx.Connection, rBuf.RemainingSlice())
		if err != nil {
			return 0, err
		}
		rBuf.Produce(int(n))
		ctx.RequestPendingBytes = max(0, ctx.RequestPendingBytes-int(n))
	}

	var i int
	for i = 0; i < len(rs); i++ {
		r := &rs[i]
		r.RemoteAddr = ctx.ClientAddress
		r.Headers = r.Headers[:0]
		r.Body = r.Body[:0]
		r.Form = r.Form[:0]
		r.Arena.Reset()

		n, err := parser.Parse(rBuf.UnconsumedString(), r)
		if err != nil {
			return i, err
		}
		if n == 0 {
			break
		}
		rBuf.Consume(n)
	}
	if (usesQ) && ((ctx.RequestPendingBytes > 0) || (i == len(rs))) {
		ctx.EventQueue.AppendEvent(event.Event{Type: event.Read, Identifier: ctx.Connection, Available: ctx.RequestPendingBytes, UserData: unsafe.Pointer(ctx)})
	}

	return i, nil
}

func ContentTypeHTML(bodies []syscall.Iovec) bool {
	if len(bodies) == 0 {
		return false
	}
	return bodies[0] == HTMLHeader
}

/* WriteResponses generates  responses and writes them on wire. Returns the number of processed responses. */
func WriteResponses(ctx *Context, ws []Response) (int, error) {
	dateBuf := ctx.DateRFC822

	for i := 0; i < len(ws); i++ {
		w := &ws[i]

		ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("HTTP/1.1"), syscall.Iovec(" "), syscall.Iovec(Status2String[w.StatusCode]), syscall.Iovec(" "), syscall.Iovec(Status2Reason[w.StatusCode]), syscall.Iovec("\r\n"))

		if !w.Headers.OmitDate {
			if dateBuf == nil {
				dateBuf := make([]byte, 31)

				var tp syscall.Timespec
				syscall.ClockGettime(syscall.CLOCK_REALTIME, &tp)
				time.PutTmRFC822(dateBuf, time.ToTm(int(tp.Sec)))
			}

			ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("Date: "), syscall.IovecForByteSlice(dateBuf), syscall.Iovec("\r\n"))
		}

		if !w.Headers.OmitServer {
			ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("Server: gofa/http\r\n"))
		}

		if !w.Headers.OmitContentType {
			if ContentTypeHTML(w.Bodies) {
				ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("Content-Type: text/html; charset=\"UTF-8\"\r\n"))
			} else {
				ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("Content-Type: text/plain; charset=\"UTF-8\"\r\n"))
			}
		}

		if !w.Headers.OmitContentLength {
			var length int
			for i := 0; i < len(w.Bodies); i++ {
				length += int(len(w.Bodies[i]))
			}

			lengthBuf := w.Arena.NewSlice(20)
			n := slices.PutInt(lengthBuf, length)

			ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("Content-Length: "), syscall.IovecForByteSlice(lengthBuf[:n]), syscall.Iovec("\r\n"))
		}

		ctx.ResponseIovs = append(ctx.ResponseIovs, w.Headers.Values...)
		ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("\r\n"))
		ctx.ResponseIovs = append(ctx.ResponseIovs, w.Bodies...)

		w.StatusCode = StatusOK
		w.Headers.Values = w.Headers.Values[:0]
		w.Headers.OmitDate = false
		w.Headers.OmitServer = false
		w.Headers.OmitContentType = false
		w.Headers.OmitContentLength = false
		w.Bodies = w.Bodies[:0]
		w.Arena.Reset()
	}

	/* IOV_MAX is 1024, so F**CK ME for not sending large pipelines with one syscall!!! */
	for len(ctx.ResponseIovs[ctx.ResponsePos:]) > 0 {
		end := min(len(ctx.ResponseIovs[ctx.ResponsePos:]), syscall.IOV_MAX)
		n, err := syscall.Writev(ctx.Connection, ctx.ResponseIovs[ctx.ResponsePos:ctx.ResponsePos+end])
		if err != nil {
			return 0, err
		}

		prevPos := ctx.ResponsePos
		for (ctx.ResponsePos < len(ctx.ResponseIovs)) && (n >= int64(len(ctx.ResponseIovs[ctx.ResponsePos]))) {
			n -= int64(len(ctx.ResponseIovs[ctx.ResponsePos]))
			ctx.ResponsePos++
		}
		if ctx.ResponsePos == len(ctx.ResponseIovs) {
			ctx.ResponseIovs = ctx.ResponseIovs[:0]
			ctx.ResponsePos = 0
		} else if ctx.ResponsePos-prevPos < end {
			log.Panicf("Written %d iovs out of %d", ctx.ResponsePos-prevPos, end)
			ctx.ResponseIovs[ctx.ResponsePos] = ctx.ResponseIovs[ctx.ResponsePos][n:]
			break

			/* TODO(anton2920): as an option, gather buffers manually into some local arena. */
		}
	}

	return len(ws), nil
}

func Close(ctx *Context) error {
	c := ctx.Connection
	FreeContext(ctx)
	return syscall.Close(c)
}
