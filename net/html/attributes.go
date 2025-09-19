package html

import (
	"github.com/anton2920/gofa/alloc"
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/trace"
)

type Attributes struct {
	Class string

	Accept  string
	Action  string
	Enctype string
	Href    string
	ID      string
	Method  string
	Name    string
	Src     string
	Type    string
	Value   string

	Alt         string
	Placeholder string

	Cols      int
	Max       int
	MaxLength int
	Min       int
	MinLength int
	Rows      int

	Checked        bool
	Disabled       bool
	FormNoValidate bool
	Multiple       bool
	Readonly       bool
	Required       bool
	Selected       bool
}

func DisplayBoolAttribute(h *HTML, attr string, value bool) {
	if value {
		h.SP()
		h.String(attr)
	}
}

func DisplayIntAttribute(h *HTML, attr string, value int) {
	if value > 0 {
		h.SP()
		h.String(attr)
		h.String(`="`)
		h.Int(value)
		h.String(`"`)
	}
}

func DisplayStringAttribute(h *HTML, attr string, value string) {
	if len(value) > 0 {
		h.SP()
		h.String(attr)
		h.String(`="`)
		h.String(value)
		h.String(`"`)
	}
}

func DisplayLStringAttribute(h *HTML, attr string, value string) {
	DisplayStringAttribute(h, attr, h.L(value))
}

func ReplaceBool(r *bool, b bool) {
	if b {
		*r = b
	}
}

func ReplaceInt(r *int, n int) {
	if n > 0 {
		*r = n
	}
}

func ReplaceString(r *string, s string) {
	if len(s) > 0 {
		*r = s
	}
}

func MergeString(arena *alloc.Arena, r *string, s string) {
	if len(s) > 0 {
		var n int

		buffer := arena.NewSlice(len(*r) + len(" ") + len(s))
		if len(*r) > 0 {
			n += copy(buffer[n:], *r)
			n += copy(buffer[n:], " ")
		}
		n += copy(buffer[n:], s)

		*r = bytes.AsString(buffer[:n])
	}
}

func (h *HTML) MergeAttributes(attrs ...Attributes) Attributes {
	t := trace.Begin("")

	var result Attributes

	if len(attrs) == 0 {
		trace.End(t)
		return result
	} else if len(attrs) == 1 {
		trace.End(t)
		return attrs[0]
	}

	for i := 0; i < len(attrs); i++ {
		attr := &attrs[i]

		MergeString(&h.Arena, &result.Class, attr.Class)

		ReplaceString(&result.Accept, attr.Accept)
		ReplaceString(&result.Action, attr.Action)
		ReplaceString(&result.Enctype, attr.Enctype)
		ReplaceString(&result.Href, attr.Href)
		ReplaceString(&result.ID, attr.ID)
		ReplaceString(&result.Method, attr.Method)
		ReplaceString(&result.Name, attr.Name)
		ReplaceString(&result.Src, attr.Src)
		ReplaceString(&result.Type, attr.Type)
		ReplaceString(&result.Value, attr.Value)

		ReplaceString(&result.Alt, attr.Alt)
		ReplaceString(&result.Placeholder, attr.Placeholder)

		ReplaceInt(&result.Cols, attr.Cols)
		ReplaceInt(&result.Max, attr.Max)
		ReplaceInt(&result.MaxLength, attr.MaxLength)
		ReplaceInt(&result.Min, attr.Min)
		ReplaceInt(&result.MinLength, attr.MinLength)
		ReplaceInt(&result.Rows, attr.Rows)

		ReplaceBool(&result.Checked, attr.Checked)
		ReplaceBool(&result.Disabled, attr.Disabled)
		ReplaceBool(&result.FormNoValidate, attr.FormNoValidate)
		ReplaceBool(&result.Multiple, attr.Multiple)
		ReplaceBool(&result.Readonly, attr.Readonly)
		ReplaceBool(&result.Required, attr.Required)
		ReplaceBool(&result.Selected, attr.Selected)
	}

	trace.End(t)
	return result
}

func (h *HTML) PrependAttributes(sys Attributes, user []Attributes) Attributes {
	return h.MergeAttributes(sys, h.MergeAttributes(user...))
}

func (h *HTML) AppendAttributes(user []Attributes, sys Attributes) Attributes {
	return h.MergeAttributes(append(user, sys)...)
}

func Class(class string) Attributes {
	return Attributes{Class: class}
}

func Enctype(enctype string) Attributes {
	return Attributes{Enctype: enctype}
}

func FormNoValidate() Attributes {
	return Attributes{FormNoValidate: true}
}

func MaxLength(n int) Attributes {
	return Attributes{MaxLength: n}
}

func MinLength(n int) Attributes {
	return Attributes{MinLength: n}
}

func Name(name string) Attributes {
	return Attributes{Name: name}
}

func Required() Attributes {
	return Attributes{Required: true}
}

func Value(value string) Attributes {
	return Attributes{Value: value}
}
