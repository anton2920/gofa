package multipart

import (
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace/trace_"
)

func ParseFormData(ctx *context.Context, contentType string, vs *url.Values, files *Files, body []byte) bool {
	t := trace_.Begin("")

	if !strings.StartsWith(contentType, "multipart/form-data") {
		ctx.NewError().S("expected 'multipart/form-data' Content-Type, got ").Q(contentType)
		trace_.End(t)
		return false
	}

	key, boundary, ok := strings.Cut(contentType, "=")
	if (!ok) || (len(key)-len("boundary") <= 0) || (key[len(key)-len("boundary"):] != "boundary") {
		ctx.NewError().S("expected boundary in Content-Type, got '").S(key[len(key)-len("boundary"):]).S(":").S(boundary).S("'")
		trace_.End(t)
		return false
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
			ctx.NewError().S("expected new line after boundary")
			trace_.End(t)
			return false
		}
		if strings.Trim(form[pos:pos+lineEnd], "-") != strings.Trim(boundary, "-") {
			ctx.NewError().S("expected boundary got ").Q(form[pos : pos+lineEnd])
			trace_.End(t)
			return false
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
				ctx.NewError().S("expected new line after header")
				trace_.End(t)
				return false
			} else if lineEnd == 0 {
				pos += len("\r\n")
				break
			}

			header := form[pos : pos+lineEnd]

			key, value, ok := strings.Cut(header, ":")
			if !ok {
				ctx.NewError().S("invalid header")
				trace_.End(t)
				return false
			}
			value = strings.TrimSpace(value)

			switch key {
			case "Content-Disposition":
				if !strings.StartsWith(value, "form-data;") {
					ctx.NewError().S("expected 'form-data', got ").Q(value)
					trace_.End(t)
					return false
				}

				leftover := value[len("form-data;"):]
				for len(leftover) > 0 {
					var pair string

					pair, leftover, _ = strings.Cut(leftover, ";")
					if len(pair) == 0 {
						ctx.NewError().S("expected header value, got nothing")
						trace_.End(t)
						return false
					}

					key, value, ok := strings.Cut(pair, "=")
					if !ok {
						ctx.NewError().S("expected key=value, got ").Q(pair)
						trace_.End(t)
						return false
					}
					value = strings.Trim(value, `"`)

					switch strings.TrimSpace(key) {
					case "name":
						name = value
					case "filename":
						filename = value
						isFile = true
					}
				}
			case "Content-Type":
				contentType = strings.Trim(value, `"`)
			}

			pos += len(header) + len("\r\n")
		}

		/* Parsing value. */
		nextBoundary := strings.FindSubstring(form[pos:], boundary)
		if nextBoundary == -1 {
			ctx.NewError().S("expected boundary after value")
			trace_.End(t)
			return false
		}
		lineEnd = strings.FindCharReverse(form[pos:pos+nextBoundary], '\r')
		if lineEnd == -1 {
			ctx.NewError().S("expected new line after value")
			trace_.End(t)
			return false
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

	trace_.End(t)
	return true
}
