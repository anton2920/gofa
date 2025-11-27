package http

import (
	"strconv"

	"github.com/anton2920/gofa/trace"
)

type Headers struct {
	Keys   []string
	Values [][]string
}

func (hs *Headers) Add(key string, value string) {
	t := trace.Begin("")

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			hs.Values[i] = append(hs.Values[i], value)
			trace.End(t)
			return
		}
	}
	hs.Keys = append(hs.Keys, key)

	if len(hs.Values) == cap(hs.Values) {
		hs.Values = append(hs.Values, []string{value})
		trace.End(t)
		return
	}
	n := len(hs.Values)
	hs.Values = hs.Values[:n+1]
	hs.Values[n] = hs.Values[n][:0]
	hs.Values[n] = append(hs.Values[n], value)

	trace.End(t)
}

func (hs *Headers) Del(key string) {
	t := trace.Begin("")

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			copy(hs.Keys[i:], hs.Keys[i+1:])
			copy(hs.Values[i:], hs.Values[i+1:])
			hs.Keys = hs.Keys[:len(hs.Keys)-1]
			hs.Values = hs.Values[:len(hs.Values)-1]
			break
		}
	}

	trace.End(t)
}

func (hs *Headers) Get(key string) string {
	t := trace.Begin("")

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			trace.End(t)
			return hs.Values[i][0]
		}
	}

	trace.End(t)
	return ""
}

func (hs *Headers) GetInt(key string) (int, error) {
	t := trace.Begin("")

	n, err := strconv.Atoi(hs.Get(key))

	trace.End(t)
	return n, err
}

func (hs *Headers) GetMany(key string) []string {
	t := trace.Begin("")

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			trace.End(t)
			return hs.Values[i]
		}
	}

	trace.End(t)
	return nil
}

func (hs *Headers) Has(key string) bool {
	t := trace.Begin("")

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			trace.End(t)
			return true
		}
	}

	trace.End(t)
	return false
}

func (hs *Headers) Reset() {
	hs.Keys = hs.Keys[:0]
	hs.Values = hs.Values[:0]
}

func (hs *Headers) Set(key string, value string) {
	t := trace.Begin("")

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			hs.Values[i] = hs.Values[i][:0]
			hs.Values[i] = append(hs.Values[i], value)
			trace.End(t)
			return
		}
	}
	hs.Keys = append(hs.Keys, key)

	if len(hs.Values) == cap(hs.Values) {
		hs.Values = append(hs.Values, []string{value})
		trace.End(t)
		return
	}
	n := len(hs.Values)
	hs.Values = hs.Values[:n+1]
	hs.Values[n] = hs.Values[n][:0]
	hs.Values[n] = append(hs.Values[n], value)

	trace.End(t)
}
