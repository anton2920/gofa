package http

import "github.com/anton2920/gofa/prof"

type Headers struct {
	Keys   []string
	Values [][]string
}

//go:nosplit
func (hs *Headers) Add(key string, value string) {
	p := prof.Begin("")

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			hs.Values[i] = append(hs.Values[i], value)
			prof.End(p)
			return
		}
	}
	hs.Keys = append(hs.Keys, key)

	if len(hs.Values) == cap(hs.Values) {
		hs.Values = append(hs.Values, []string{value})
		prof.End(p)
		return
	}
	n := len(hs.Values)
	hs.Values = hs.Values[:n+1]
	hs.Values[n] = hs.Values[n][:0]
	hs.Values[n] = append(hs.Values[n], value)

	prof.End(p)
}

//go:nosplit
func (hs *Headers) Get(key string) string {
	p := prof.Begin("")

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			prof.End(p)
			return hs.Values[i][0]
		}
	}

	prof.End(p)
	return ""
}

//go:nosplit
func (hs *Headers) GetMany(key string) []string {
	p := prof.Begin("")

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			prof.End(p)
			return hs.Values[i]
		}
	}

	prof.End(p)
	return nil
}

//go:nosplit
func (hs *Headers) Has(key string) bool {
	p := prof.Begin("")

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			prof.End(p)
			return true
		}
	}

	prof.End(p)
	return false
}

//go:nosplit
func (hs *Headers) Reset() {
	p := prof.Begin("")

	hs.Keys = hs.Keys[:0]
	hs.Values = hs.Values[:0]

	prof.End(p)
}

//go:nosplit
func (hs *Headers) Set(key string, value string) {
	p := prof.Begin("")

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			hs.Values[i] = hs.Values[i][:0]
			hs.Values[i] = append(hs.Values[i], value)
			prof.End(p)
			return
		}
	}
	hs.Keys = append(hs.Keys, key)

	if len(hs.Values) == cap(hs.Values) {
		hs.Values = append(hs.Values, []string{value})
		prof.End(p)
		return
	}
	n := len(hs.Values)
	hs.Values = hs.Values[:n+1]
	hs.Values[n] = hs.Values[n][:0]
	hs.Values[n] = append(hs.Values[n], value)

	prof.End(p)
}
