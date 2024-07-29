package http

import "github.com/anton2920/gofa/trace"

type Headers struct {
	Keys   []string
	Values [][]string
}

func (hs *Headers) Add(key string, value string) {
	defer trace.End(trace.Start(""))

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			hs.Values[i] = append(hs.Values[i], value)
			return
		}
	}
	hs.Keys = append(hs.Keys, key)

	if len(hs.Values) == cap(hs.Values) {
		hs.Values = append(hs.Values, []string{value})
		return
	}
	n := len(hs.Values)
	hs.Values = hs.Values[:n+1]
	hs.Values[n] = hs.Values[n][:0]
	hs.Values[n] = append(hs.Values[n], value)
}

func (hs *Headers) Get(key string) string {
	defer trace.End(trace.Start(""))

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			return hs.Values[i][0]
		}
	}

	return ""
}

func (hs *Headers) GetMany(key string) []string {
	defer trace.End(trace.Start(""))

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			return hs.Values[i]
		}
	}

	return nil
}

func (hs *Headers) Has(key string) bool {
	defer trace.End(trace.Start(""))

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			return true
		}
	}

	return false
}

func (hs *Headers) Reset() {
	defer trace.End(trace.Start(""))

	hs.Keys = hs.Keys[:0]
	hs.Values = hs.Values[:0]
}

func (hs *Headers) Set(key string, value string) {
	defer trace.End(trace.Start(""))

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			hs.Values[i] = hs.Values[i][:0]
			hs.Values[i] = append(hs.Values[i], value)
			return
		}
	}
	hs.Keys = append(hs.Keys, key)

	if len(hs.Values) == cap(hs.Values) {
		hs.Values = append(hs.Values, []string{value})
		return
	}
	n := len(hs.Values)
	hs.Values = hs.Values[:n+1]
	hs.Values[n] = hs.Values[n][:0]
	hs.Values[n] = append(hs.Values[n], value)
}
