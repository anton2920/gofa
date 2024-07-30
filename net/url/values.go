package url

import (
	"strconv"
	"unsafe"

	"github.com/anton2920/gofa/arena"
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/prof"
	"github.com/anton2920/gofa/slices"
)

type Values struct {
	Arena arena.Arena

	Keys   []string
	Values [][]string
}

func (vs *Values) Add(key string, value string) {
	defer prof.End(prof.Begin(""))

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			vs.Values[i] = append(vs.Values[i], value)
			return
		}
	}
	vs.Keys = append(vs.Keys, key)

	if len(vs.Values) == cap(vs.Values) {
		vs.Values = append(vs.Values, []string{value})
		return
	}
	n := len(vs.Values)
	vs.Values = vs.Values[:n+1]
	vs.Values[n] = vs.Values[n][:0]
	vs.Values[n] = append(vs.Values[n], value)
}

//go:nosplit
func (vs *Values) Get(key string) string {
	defer prof.End(prof.Begin(""))

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			return vs.Values[i][0]
		}
	}

	return ""
}

//go:nosplit
func (vs Values) GetInt(key string) (int, error) {
	defer prof.End(prof.Begin(""))

	return strconv.Atoi(vs.Get(key))
}

//go:nosplit
func (vs Values) GetID(key string) (database.ID, error) {
	defer prof.End(prof.Begin(""))

	id, err := strconv.Atoi(vs.Get(key))
	if err != nil {
		return -1, err
	}
	if (id < database.MinValidID) || (id > database.MaxValidID) {
		return -1, errors.New("ID out of range")
	}
	return database.ID(id), nil
}

//go:nosplit
func (vs *Values) GetMany(key string) []string {
	defer prof.End(prof.Begin(""))

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			return vs.Values[i]
		}
	}

	return nil
}

//go:nosplit
func (vs *Values) Has(key string) bool {
	defer prof.End(prof.Begin(""))

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			return true
		}
	}

	return false
}

//go:nosplit
func (vs *Values) Reset() {
	defer prof.End(prof.Begin(""))

	vs.Keys = vs.Keys[:0]
	vs.Values = vs.Values[:0]
	vs.Arena.Reset()
}

//go:nosplit
func (vs *Values) Set(key string, value string) {
	defer prof.End(prof.Begin(""))

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			vs.Values[i] = vs.Values[i][:0]
			vs.Values[i] = append(vs.Values[i], value)
			return
		}
	}
	vs.Keys = append(vs.Keys, key)

	if len(vs.Values) == cap(vs.Values) {
		vs.Values = append(vs.Values, []string{value})
		return
	}
	n := len(vs.Values)
	vs.Values = vs.Values[:n+1]
	vs.Values[n] = vs.Values[n][:0]
	vs.Values[n] = append(vs.Values[n], value)
}

//go:nosplit
func (vs *Values) SetInt(key string, value int) {
	defer prof.End(prof.Begin(""))

	buffer := vs.Arena.NewSlice(20)
	n := slices.PutInt(buffer, value)
	vs.Set(key, unsafe.String(&buffer[0], n))
}
