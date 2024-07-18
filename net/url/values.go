package url

import (
	"errors"
	"strconv"
	"unsafe"

	"github.com/anton2920/gofa/arena"
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/slices"
)

type Values struct {
	Arena arena.Arena

	Keys   []string
	Values [][]string
}

func (vs *Values) Add(key string, value string) {
	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			vs.Values[i] = append(vs.Values[i], value)
			return
		}
	}

	vs.Keys = append(vs.Keys, key)
	vs.Values = append(vs.Values, []string{value})
}

func (vs *Values) Get(key string) string {
	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			return vs.Values[i][0]
		}
	}

	return ""
}

func (vs Values) GetInt(key string) (int, error) {
	return strconv.Atoi(vs.Get(key))
}

func (vs Values) GetID(key string) (database.ID, error) {
	id, err := strconv.Atoi(vs.Get(key))
	if err != nil {
		return -1, err
	}
	if (id < database.MinValidID) || (id > database.MaxValidID) {
		return -1, errors.New("ID out of range")
	}
	return database.ID(id), nil
}

func (vs *Values) GetMany(key string) []string {
	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			return vs.Values[i]
		}
	}

	return nil
}

func (vs *Values) Has(key string) bool {
	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			return true
		}
	}

	return false
}

func (vs *Values) Reset() {
	vs.Keys = vs.Keys[:0]

	for i := 0; i < len(vs.Values); i++ {
		vs.Values[i] = vs.Values[i][:0]
	}
	vs.Values = vs.Values[:0]

	vs.Arena.Reset()
}

func (vs *Values) Set(key string, value string) {
	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			vs.Values[i] = vs.Values[i][:0]
			vs.Values[i] = append(vs.Values[i], value)
			return
		}
	}

	vs.Keys = append(vs.Keys, key)
	vs.Values = append(vs.Values, []string{value})
}

func (vs *Values) SetInt(key string, value int) {
	buffer := vs.Arena.NewSlice(20)
	n := slices.PutInt(buffer, value)
	vs.Set(key, unsafe.String(&buffer[0], n))
}
