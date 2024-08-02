package url

import (
	"strconv"
	"unsafe"

	"github.com/anton2920/gofa/alloc"
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/prof"
	"github.com/anton2920/gofa/slices"
)

type Values struct {
	Arena alloc.Arena

	Keys   []string
	Values [][]string
}

func (vs *Values) Add(key string, value string) {
	p := prof.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			vs.Values[i] = append(vs.Values[i], value)
			prof.End(p)
			return
		}
	}
	vs.Keys = append(vs.Keys, key)

	if len(vs.Values) == cap(vs.Values) {
		vs.Values = append(vs.Values, []string{value})
		prof.End(p)
		return
	}
	n := len(vs.Values)
	vs.Values = vs.Values[:n+1]
	vs.Values[n] = vs.Values[n][:0]
	vs.Values[n] = append(vs.Values[n], value)

	prof.End(p)
}

func (vs *Values) Get(key string) string {
	p := prof.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			prof.End(p)
			return vs.Values[i][0]
		}
	}

	prof.End(p)
	return ""
}

func (vs Values) GetInt(key string) (int, error) {
	p := prof.Begin("")

	n, err := strconv.Atoi(vs.Get(key))

	prof.End(p)
	return n, err
}

func (vs Values) GetID(key string) (database.ID, error) {
	p := prof.Begin("")

	id, err := strconv.Atoi(vs.Get(key))
	if err != nil {
		prof.End(p)
		return -1, err
	}
	if (id < database.MinValidID) || (id > database.MaxValidID) {
		prof.End(p)
		return -1, errors.New("ID out of range")
	}

	prof.End(p)
	return database.ID(id), nil
}

func (vs *Values) GetMany(key string) []string {
	p := prof.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			prof.End(p)
			return vs.Values[i]
		}
	}

	prof.End(p)
	return nil
}

func (vs *Values) Has(key string) bool {
	p := prof.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			prof.End(p)
			return true
		}
	}

	prof.End(p)
	return false
}

func (vs *Values) Reset() {
	vs.Keys = vs.Keys[:0]
	vs.Values = vs.Values[:0]
	vs.Arena.Reset()
}

func (vs *Values) Set(key string, value string) {
	p := prof.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			vs.Values[i] = vs.Values[i][:0]
			vs.Values[i] = append(vs.Values[i], value)
			prof.End(p)
			return
		}
	}
	vs.Keys = append(vs.Keys, key)

	if len(vs.Values) == cap(vs.Values) {
		vs.Values = append(vs.Values, []string{value})
		prof.End(p)
		return
	}
	n := len(vs.Values)
	vs.Values = vs.Values[:n+1]
	vs.Values[n] = vs.Values[n][:0]
	vs.Values[n] = append(vs.Values[n], value)

	prof.End(p)
}

func (vs *Values) SetInt(key string, value int) {
	p := prof.Begin("")

	buffer := vs.Arena.NewSlice(20)
	n := slices.PutInt(buffer, value)
	vs.Set(key, unsafe.String(&buffer[0], n))

	prof.End(p)
}
