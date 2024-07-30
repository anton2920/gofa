package http

import "github.com/anton2920/gofa/prof"

type Headers struct {
	Keys   []string
	Values [][]string
}

//go:nosplit
func (hs *Headers) Add(key string, value string) {
	defer prof.End(prof.Begin(""))

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

//go:nosplit
func (hs *Headers) Get(key string) string {
	defer prof.End(prof.Begin(""))

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			return hs.Values[i][0]
		}
	}

	return ""
}

//go:nosplit
func (hs *Headers) GetMany(key string) []string {
	defer prof.End(prof.Begin(""))

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			return hs.Values[i]
		}
	}

	return nil
}

//go:nosplit
func (hs *Headers) Has(key string) bool {
	defer prof.End(prof.Begin(""))

	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			return true
		}
	}

	return false
}

//go:nosplit
func (hs *Headers) Reset() {
	defer prof.End(prof.Begin(""))

	hs.Keys = hs.Keys[:0]
	hs.Values = hs.Values[:0]
}

//go:nosplit
func (hs *Headers) Set(key string, value string) {
	defer prof.End(prof.Begin(""))

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
