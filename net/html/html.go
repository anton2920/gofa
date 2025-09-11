package html

import (
	"github.com/anton2920/gofa/alloc"
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/l10n"
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/trace"
)

type ApplyType int32

type Options struct {
	ApplyType

	ID    string
	Class string

	Type        string
	Name        string
	Value       string
	Placeholder string

	Min int
	Max int

	Cols int
	Rows int

	MinLength int
	MaxLength int

	Disabled       bool
	Readonly       bool
	Required       bool
	FormNoValidate bool
}

const (
	ApplyTypeMerge = ApplyType(iota)
	ApplyTypeReplace
	ApplyTypeCount
)

func ProcessOptions(arena *alloc.Arena, dflt Options, user Options) Options {
	t := trace.Begin("")

	if user.ApplyType == ApplyTypeReplace {
		trace.End(t)
		return user
	}

	result := dflt

	if len(user.ID) > 0 {
		result.ID = user.ID
	}
	if len(user.Class) > 0 {
		var n int

		buffer := arena.NewSlice(len(dflt.Class) + len(" ") + len(user.Class))
		n += copy(buffer[n:], dflt.Class)
		n += copy(buffer[n:], " ")
		n += copy(buffer[n:], user.Class)

		result.Class = bytes.AsString(buffer)
	}

	if len(user.Type) > 0 {
		result.Type = user.Type
	}
	if len(user.Name) > 0 {
		result.Name = user.Name
	}
	if len(user.Value) > 0 {
		result.Value = user.Value
	}
	if len(user.Placeholder) > 0 {
		result.Placeholder = user.Placeholder
	}

	if user.Min > 0 {
		result.Min = user.Min
	}
	if user.Max > 0 {
		result.Max = user.Max
	}

	if user.MinLength > 0 {
		result.MinLength = user.MinLength
	}
	if user.MaxLength > 0 {
		result.MaxLength = user.MaxLength
	}

	if user.Disabled {
		result.Disabled = user.Disabled
	}
	if user.Readonly {
		result.Readonly = user.Readonly
	}
	if user.Required {
		result.Required = user.Required
	}

	trace.End(t)
	return result
}

func TagBegin(w *http.Response, tag string, opts Options) {
	t := trace.Begin("")

	opts = ProcessOptions(&w.Arena, Options{}, opts)

	w.WriteString(`<`)
	w.WriteString(tag)

	if len(opts.ID) > 0 {
		w.WriteString(` id="`)
		w.WriteString(opts.ID)
		w.WriteString(`"`)
	}
	if len(opts.Class) > 0 {
		w.WriteString(` class=`)
		w.WriteString(opts.Class)
		w.WriteString(`"`)
	}

	if len(opts.Type) > 0 {
		w.WriteString(` type=`)
		w.WriteString(opts.Type)
		w.WriteString(`"`)
	}
	if len(opts.Name) > 0 {
		w.WriteString(` name=`)
		w.WriteString(opts.Name)
		w.WriteString(`"`)
	}
	if len(opts.Value) > 0 {
		w.WriteString(` value=`)
		w.WriteString(opts.Value)
		w.WriteString(`"`)
	}
	if len(opts.Placeholder) > 0 {
		w.WriteString(` placeholder=`)
		w.WriteString(opts.Placeholder)
		w.WriteString(`"`)
	}

	if opts.Min > 0 {
		w.WriteString(` min=`)
		w.WriteInt(opts.Min)
		w.WriteString(`"`)
	}
	if opts.Max > 0 {
		w.WriteString(` max=`)
		w.WriteInt(opts.Max)
		w.WriteString(`"`)
	}

	if opts.Cols > 0 {
		w.WriteString(` cols=`)
		w.WriteInt(opts.Cols)
		w.WriteString(`"`)
	}
	if opts.Rows > 0 {
		w.WriteString(` rows=`)
		w.WriteInt(opts.Rows)
		w.WriteString(`"`)
	}

	if opts.MinLength > 0 {
		w.WriteString(` minlength=`)
		w.WriteInt(opts.MinLength)
		w.WriteString(`"`)
	}
	if opts.MaxLength > 0 {
		w.WriteString(` maxlength=`)
		w.WriteInt(opts.MaxLength)
		w.WriteString(`"`)
	}

	if opts.Disabled {
		w.WriteString(` disabled`)
	}
	if opts.Readonly {
		w.WriteString(` readonly`)
	}
	if opts.Required {
		w.WriteString(` required`)
	}
	if opts.FormNoValidate {
		w.WriteString(` formnovalidate`)
	}

	w.WriteString(`>`)

	trace.End(t)
}

func TagEnd(w *http.Response, tag string) {
	t := trace.Begin("")

	w.WriteString(`</`)
	w.WriteString(tag)
	w.WriteString(`>`)

	trace.End(t)
}

func Begin(w *http.Response, l l10n.Language) {
	t := trace.Begin("")

	w.WriteString(`<!DOCTYPE html`)
	w.WriteString(`<html lang="`)
	w.WriteString(l10n.Language2HTMLLang[l])
	w.WriteString(`">`)

	trace.End(t)
}

func End(w *http.Response) {
	t := trace.Begin("")

	w.WriteString(`</html>`)

	trace.End(t)
}
