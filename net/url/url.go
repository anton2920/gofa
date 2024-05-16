package url

import "strconv"

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

/* CharToByte returns ASCII-decoded character. For example, 'A' yields '\x0A'. */
func CharToByte(c byte) (byte, bool) {
	if c >= '0' && c <= '9' {
		return c - '0', true
	} else if c >= 'A' && c <= 'F' {
		return 10 + c - 'A', true
	} else {
		return '\x00', false
	}
}

func QueryDecode(decoded []byte, encoded string) (int, bool) {
	var hi, lo byte
	var ok bool
	var n int

	for i := 0; i < len(encoded); i++ {
		if encoded[i] == '%' {
			hi = encoded[i+1]
			hi, ok = CharToByte(hi)
			if !ok {
				return 0, false
			}

			lo = encoded[i+2]
			lo, ok = CharToByte(lo)
			if !ok {
				return 0, false
			}

			decoded[n] = byte(hi<<4 | lo)
			i += 2
		} else if encoded[i] == '+' {
			decoded[n] = ' '
		} else {
			decoded[n] = encoded[i]
		}
		n++
	}
	return n, true
}
