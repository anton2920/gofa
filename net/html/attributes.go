package html

import (
	"github.com/anton2920/gofa/alloc"
	"github.com/anton2920/gofa/bools"
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/strings"
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
	Rel     string
	Src     string
	Style   string
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

		strings.Replace(&result.Accept, attr.Accept)
		strings.Replace(&result.Action, attr.Action)
		strings.Replace(&result.Enctype, attr.Enctype)
		strings.Replace(&result.Href, attr.Href)
		strings.Replace(&result.ID, attr.ID)
		strings.Replace(&result.Method, attr.Method)
		strings.Replace(&result.Name, attr.Name)
		strings.Replace(&result.Rel, attr.Rel)
		strings.Replace(&result.Src, attr.Src)
		strings.Replace(&result.Style, attr.Style)
		strings.Replace(&result.Type, attr.Type)
		strings.Replace(&result.Value, attr.Value)

		strings.Replace(&result.Alt, attr.Alt)
		strings.Replace(&result.Placeholder, attr.Placeholder)

		ints.Replace(&result.Cols, attr.Cols)
		ints.Replace(&result.Max, attr.Max)
		ints.Replace(&result.MaxLength, attr.MaxLength)
		ints.Replace(&result.Min, attr.Min)
		ints.Replace(&result.MinLength, attr.MinLength)
		ints.Replace(&result.Rows, attr.Rows)

		bools.Replace(&result.Checked, attr.Checked)
		bools.Replace(&result.Disabled, attr.Disabled)
		bools.Replace(&result.FormNoValidate, attr.FormNoValidate)
		bools.Replace(&result.Multiple, attr.Multiple)
		bools.Replace(&result.Readonly, attr.Readonly)
		bools.Replace(&result.Required, attr.Required)
		bools.Replace(&result.Selected, attr.Selected)
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

func Action(action string) Attributes {
	return Attributes{Action: action}
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

func Style(s string) Attributes {
	return Attributes{Style: s}
}

func Value(value string) Attributes {
	return Attributes{Value: value}
}
