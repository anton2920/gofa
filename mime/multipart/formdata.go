package multipart

import (
	"fmt"
	stdstrings "strings"

	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace"
)

func ParseFormData(contentType string, vs *url.Values, files *Files, body []byte) error {
	t := trace.Begin("")

	if !strings.StartsWith(contentType, "multipart/form-data") {
		trace.End(t)
		return fmt.Errorf("expected 'multipart/form-data' Content-Type, got %q", contentType)
	}

	key, boundary, ok := strings.Cut(contentType, "=")
	if (!ok) || (len(key)-len("boundary") <= 0) || (key[len(key)-len("boundary"):] != "boundary") {
		trace.End(t)
		return fmt.Errorf("expected boundary in Content-Type, got '%s:%s'", key[len(key)-len("boundary"):], boundary)
	}

	form := bytes.AsString(body)
	var pos int
	for {
		/* Parsing boundary. */
		lineEnd := strings.FindChar(form[pos:], '\r')
		if lineEnd == 0 {
			break
		}
		if lineEnd == -1 {
			trace.End(t)
			return errors.New("expected new line after boundary")
		}
		if stdstrings.Trim(form[pos:pos+lineEnd], "-") != stdstrings.Trim(boundary, "-") {
			trace.End(t)
			return fmt.Errorf("expected boundary got %q", form[pos:pos+lineEnd])
		}
		if form[pos+lineEnd-2:pos+lineEnd] == "--" {
			break
		}
		pos += lineEnd + len("\r\n")

		/* Parsing headers. */
		var name, filename, contentType string
		var isFile bool

		for {
			lineEnd := strings.FindChar(form[pos:], '\r')
			if lineEnd == -1 {
				trace.End(t)
				return errors.New("expected new line after header")
			} else if lineEnd == 0 {
				pos += len("\r\n")
				break
			}

			header := form[pos : pos+lineEnd]

			key, value, ok := strings.Cut(header, ":")
			if !ok {
				trace.End(t)
				return errors.New("invalid header")
			}
			value = stdstrings.TrimSpace(value)

			switch key {
			case "Content-Disposition":
				if !strings.StartsWith(value, "form-data;") {
					trace.End(t)
					return fmt.Errorf("expected 'form-data', got %q", value)
				}

				leftover := value[len("form-data;"):]
				for len(leftover) > 0 {
					var pair string

					pair, leftover, _ = strings.Cut(leftover, ";")
					if len(pair) == 0 {
						trace.End(t)
						return errors.New("expected header value, got nothing")
					}

					key, value, ok := strings.Cut(pair, "=")
					if !ok {
						trace.End(t)
						return fmt.Errorf("expected key=value, got %q", pair)
					}
					value = stdstrings.Trim(value, `"`)

					switch stdstrings.TrimSpace(key) {
					case "name":
						name = value
					case "filename":
						filename = value
						isFile = true
					}
				}
			case "Content-Type":
				contentType = stdstrings.Trim(value, `"`)
			}

			pos += len(header) + len("\r\n")
		}

		/* Parsing value. */
		nextBoundary := strings.FindSubstring(form[pos:], boundary)
		if nextBoundary == -1 {
			trace.End(t)
			return errors.New("expected boundary after value")
		}
		lineEnd = strings.FindCharReverse(form[pos:pos+nextBoundary], '\r')
		if lineEnd == -1 {
			trace.End(t)
			return errors.New("expected new line after value")
		}
		value := form[pos : pos+lineEnd]
		if len(name) > 0 {
			if isFile {
				files.Add(name, File{filename, contentType, strings.AsBytes(value)})
			} else {
				vs.Add(name, value)
			}
		}
		pos += lineEnd + len("\r\n")
	}

	trace.End(t)
	return nil
}
