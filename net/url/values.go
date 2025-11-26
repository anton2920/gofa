package url

import (
	"strconv"

	"github.com/anton2920/gofa/alloc"
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace"
)

type Values struct {
	Keys   []string
	Values [][]string
}

func ParseQuery(arena *alloc.Arena, vs *Values, query string) error {
	t := trace.Begin("")

	var err error

	for query != "" {
		var key string
		key, query, _ = strings.Cut(query, "&")
		if strings.FindChar(key, ';') != -1 {
			err = errors.New("invalid semicolon separator in query")
			continue
		}
		if key == "" {
			continue
		}
		key, value, _ := strings.Cut(key, "=")

		keyBuffer := arena.NewSlice(len(key))
		n, ok := QueryDecode(keyBuffer, key)
		if !ok {
			if err == nil {
				err = errors.New("invalid key")
			}
			continue
		}
		key = bytes.AsString(keyBuffer[:n])

		valueBuffer := arena.NewSlice(len(value))
		n, ok = QueryDecode(valueBuffer, value)
		if !ok {
			if err == nil {
				err = errors.New("invalid value")
			}
			continue
		}
		value = bytes.AsString(valueBuffer[:n])

		vs.Add(key, value)
	}

	trace.End(t)
	return err
}

func RemoveStringAtIndex(vs []string, i int) []string {
	if (len(vs) == 0) || (i < 0) || (i >= len(vs)) {
		return vs
	}
	if i < len(vs)-1 {
		copy(vs[i:], vs[i+1:])
	}
	return vs[:len(vs)-1]
}

func RemoveValuesAtIndex(vs [][]string, i int) [][]string {
	if (len(vs) == 0) || (i < 0) || (i >= len(vs)) {
		return vs
	}
	if i < len(vs)-1 {
		copy(vs[i:], vs[i+1:])
	}
	return vs[:len(vs)-1]
}

func (vs *Values) Add(key string, value string) {
	t := trace.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			vs.Values[i] = append(vs.Values[i], value)

			trace.End(t)
			return
		}
	}
	vs.Keys = append(vs.Keys, key)

	if len(vs.Values) == cap(vs.Values) {
		vs.Values = append(vs.Values, []string{value})

		trace.End(t)
		return
	}
	n := len(vs.Values)
	vs.Values = vs.Values[:n+1]
	vs.Values[n] = vs.Values[n][:0]
	vs.Values[n] = append(vs.Values[n], value)

	trace.End(t)
}

func (vs *Values) Del(key string) {
	t := trace.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			vs.Keys = RemoveStringAtIndex(vs.Keys, i)
			vs.Values = RemoveValuesAtIndex(vs.Values, i)
			break
		}
	}

	trace.End(t)
}

func (vs *Values) Get(key string) string {
	t := trace.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			trace.End(t)
			return vs.Values[i][0]
		}
	}

	trace.End(t)
	return ""
}

func (vs Values) GetInt(key string) (int, error) {
	n, err := vs.GetInt64(key)
	return int(n), err
}

func (vs Values) GetInt32(key string) (int32, error) {
	n, err := vs.GetInt64(key)
	return int32(n), err
}

func (vs Values) GetInt64(key string) (int64, error) {
	t := trace.Begin("")

	n, err := strconv.ParseInt(vs.Get(key), 10, 64)

	trace.End(t)
	return n, err
}

func (vs Values) GetID(key string) (database.ID, error) {
	t := trace.Begin("")

	id, err := vs.GetInt(key)
	if err != nil {
		trace.End(t)
		return 0, err
	}
	if (id < database.MinValidID) || (id > database.MaxValidID) {
		trace.End(t)
		return 0, errors.New("ID out of range")
	}

	trace.End(t)
	return database.ID(id), nil
}

func (vs *Values) GetMany(key string) []string {
	t := trace.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			trace.End(t)
			return vs.Values[i]
		}
	}

	trace.End(t)
	return nil
}

func (vs *Values) Has(key string) bool {
	t := trace.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			trace.End(t)
			return true
		}
	}

	trace.End(t)
	return false
}

func (vs *Values) HasID(id database.ID) bool {
	t := trace.Begin("")

	buffer := make([]byte, ints.Bufsize)
	n := slices.PutInt(buffer, int(id))
	has := vs.Has(bytes.AsString(buffer[:n]))

	trace.End(t)
	return has
}

func (vs *Values) HasInt(value int) bool {
	t := trace.Begin("")

	buffer := make([]byte, ints.Bufsize)
	n := slices.PutInt(buffer, value)
	has := vs.Has(bytes.AsString(buffer[:n]))

	trace.End(t)
	return has
}

func (vs *Values) Reset() {
	vs.Keys = vs.Keys[:0]
	vs.Values = vs.Values[:0]
}

func (vs *Values) Set(key string, value string) {
	t := trace.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			vs.Values[i] = vs.Values[i][:0]
			vs.Values[i] = append(vs.Values[i], value)

			trace.End(t)
			return
		}
	}
	vs.Keys = append(vs.Keys, key)

	if len(vs.Values) == cap(vs.Values) {
		vs.Values = append(vs.Values, []string{value})

		trace.End(t)
		return
	}
	n := len(vs.Values)
	vs.Values = vs.Values[:n+1]
	vs.Values[n] = vs.Values[n][:0]
	vs.Values[n] = append(vs.Values[n], value)

	trace.End(t)
}

/* TODO(anton2920): remove this function altogether. */
func (vs *Values) SetInt(key string, value int) {
	t := trace.Begin("")

	buffer := make([]byte, ints.Bufsize)
	n := slices.PutInt(buffer, value)
	vs.Set(key, string(buffer[:n]))

	trace.End(t)
}

func (vs *Values) SetID(key string, id database.ID) {
	vs.SetInt(key, int(id))
}
