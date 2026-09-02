package url

import (
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace/trace_"
)

type Values struct {
	Keys   []string
	Values [][]string
}

func ParseQuery(ctx *context.Context, vs *Values, query string) bool {
	t := trace_.Begin("")

	for query != "" {
		var key string
		key, query, _ = strings.Cut(query, "&")
		if strings.FindChar(key, ';') != -1 {
			ctx.NewError().S("invalid semicolon separator in query")
			return false
		}
		if key == "" {
			continue
		}
		key, value, _ := strings.Cut(key, "=")

		keyBuffer := ctx.Arena.PushByteArray(len(key))
		n, ok := QueryDecode(keyBuffer, key)
		if !ok {
			ctx.NewError().S("invalid key")
			return false
		}
		key = bytes.AsString(keyBuffer[:n])

		valueBuffer := ctx.Arena.PushByteArray(len(value))
		n, ok = QueryDecode(valueBuffer, value)
		if !ok {
			ctx.NewError().S("invalid value")
			return false
		}
		value = bytes.AsString(valueBuffer[:n])

		vs.Add(key, value)
	}

	trace_.End(t)
	return true
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
	t := trace_.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			vs.Values[i] = append(vs.Values[i], value)

			trace_.End(t)
			return
		}
	}
	vs.Keys = append(vs.Keys, key)

	if len(vs.Values) == cap(vs.Values) {
		vs.Values = append(vs.Values, []string{value})

		trace_.End(t)
		return
	}
	n := len(vs.Values)
	vs.Values = vs.Values[:n+1]
	vs.Values[n] = vs.Values[n][:0]
	vs.Values[n] = append(vs.Values[n], value)

	trace_.End(t)
}

func (vs *Values) Del(key string) {
	t := trace_.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			vs.Keys = RemoveStringAtIndex(vs.Keys, i)
			vs.Values = RemoveValuesAtIndex(vs.Values, i)
			break
		}
	}

	trace_.End(t)
}

func (vs *Values) Get(key string) string {
	t := trace_.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			trace_.End(t)
			return vs.Values[i][0]
		}
	}

	trace_.End(t)
	return ""
}

func (vs Values) GetInt(key string) (int, error) {
	n, err := 0, error(nil) //vs.GetInt64(key)
	return int(n), err
}

/*
func (vs Values) GetInt32(key string) (int32, error) {
	n, err := vs.GetInt64(key)
	return int32(n), err
}

func (vs Values) GetInt64(key string) (int64, error) {
	t := trace_.Begin("")

	n, err := strconv.ParseInt(vs.Get(key), 10, 64)

	trace_.End(t)
	return n, err
}
*/

func (vs *Values) GetMany(key string) []string {
	t := trace_.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			trace_.End(t)
			return vs.Values[i]
		}
	}

	trace_.End(t)
	return nil
}

func (vs *Values) Has(key string) bool {
	t := trace_.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			trace_.End(t)
			return true
		}
	}

	trace_.End(t)
	return false
}

func (vs *Values) HasInt(value int) bool {
	t := trace_.Begin("")

	buffer := make([]byte, ints.Bufsize)
	n := slices.PutInt(buffer, value)
	has := vs.Has(bytes.AsString(buffer[:n]))

	trace_.End(t)
	return has
}

func (vs *Values) Reset() {
	vs.Keys = vs.Keys[:0]
	vs.Values = vs.Values[:0]
}

func (vs *Values) Set(key string, value string) {
	t := trace_.Begin("")

	for i := 0; i < len(vs.Keys); i++ {
		if key == vs.Keys[i] {
			vs.Values[i] = vs.Values[i][:0]
			vs.Values[i] = append(vs.Values[i], value)

			trace_.End(t)
			return
		}
	}
	vs.Keys = append(vs.Keys, key)

	if len(vs.Values) == cap(vs.Values) {
		vs.Values = append(vs.Values, []string{value})

		trace_.End(t)
		return
	}
	n := len(vs.Values)
	vs.Values = vs.Values[:n+1]
	vs.Values[n] = vs.Values[n][:0]
	vs.Values[n] = append(vs.Values[n], value)

	trace_.End(t)
}

/* TODO(anton2920): remove this function altogether. */
func (vs *Values) SetInt(key string, value int) {
	t := trace_.Begin("")

	buffer := make([]byte, ints.Bufsize)
	n := slices.PutInt(buffer, value)
	vs.Set(key, string(buffer[:n]))

	trace_.End(t)
}
