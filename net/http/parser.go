package http

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/strings"
)

type State int

const (
	StateRequestLine State = iota
	StateHeader
	StateBody

	StateUnknown
	StateDone
)

type RequestParser struct {
	State State
	Pos   int

	ContentLength int
}

func (rp *RequestParser) Parse(request string, r *Request) (int, error) {
	var err error

	rp.State = StateRequestLine
	rp.ContentLength = 0
	rp.Pos = 0

	for rp.State != StateDone {
		switch rp.State {
		default:
			log.Panicf("Unknown parser state %d", rp.State)
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
				return 0, fmt.Errorf("expected space after URI, found %q", request[rp.Pos:rp.Pos+lineEnd])
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
