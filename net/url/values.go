package url

import (
	"strconv"

	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/errors"
)

type URL struct {
	Path  string
	Query string
}

type Value struct {
	Key    string
	Values []string
}

type Values []Value

func (vs *Values) Add(key string, value string) {
	if vs == nil {
		*vs = append(*vs, Value{Key: key, Values: []string{value}})
		return
	}

	for i := 0; i < len(*vs); i++ {
		v := &(*vs)[i]
		if key == v.Key {
			v.Values = append(v.Values, value)
			return
		}
	}

	if len(*vs) >= cap(*vs) {
		*vs = append(*vs, Value{Key: key, Values: []string{value}})
		return
	}

	l := len(*vs)
	*vs = (*vs)[:l+1]

	v := &(*vs)[l]
	v.Key = key
	v.Values = v.Values[:0]
	v.Values = append(v.Values, value)
}

func (vs Values) Get(key string) string {
	for i := 0; i < len(vs); i++ {
		if key == vs[i].Key {
			values := vs[i].Values
			if len(values) == 0 {
				return ""
			}
			return values[0]
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

func (vs Values) GetMany(key string) []string {
	for i := 0; i < len(vs); i++ {
		if key == vs[i].Key {
			return vs[i].Values
		}
	}
	return nil
}

func (vs *Values) Set(key string, value string) {
	if vs == nil {
		*vs = append(*vs, Value{Key: key, Values: []string{value}})
		return
	}

	for i := 0; i < len(*vs); i++ {
		v := &(*vs)[i]
		if key == v.Key {
			v.Values = v.Values[:0]
			v.Values = append(v.Values, value)
			return
		}
	}

	if len(*vs) >= cap(*vs) {
		*vs = append(*vs, Value{Key: key, Values: []string{value}})
		return
	}

	l := len(*vs)
	*vs = (*vs)[:l+1]

	v := &(*vs)[l]
	v.Key = key
	v.Values = v.Values[:0]
	v.Values = append(v.Values, value)
}

func (vs *Values) SetInt(key string, value int) {
	/* TODO(anton2920): change to arena allocation. */
	vs.Set(key, strconv.Itoa(value))
}
