package http

type Headers struct {
	Keys   []string
	Values [][]string
}

func (hs *Headers) Add(key string, value string) {
	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			hs.Values[i] = append(hs.Values[i], value)
			return
		}
	}

	hs.Keys = append(hs.Keys, key)
	hs.Values = append(hs.Values, []string{value})
}

func (hs *Headers) Get(key string) string {
	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			return hs.Values[i][0]
		}
	}

	return ""
}

func (hs *Headers) GetMany(key string) []string {
	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			return hs.Values[i]
		}
	}

	return nil
}

func (hs *Headers) Has(key string) bool {
	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			return true
		}
	}

	return false
}

func (hs *Headers) Reset() {
	hs.Keys = hs.Keys[:0]

	for i := 0; i < len(hs.Values); i++ {
		hs.Values[i] = hs.Values[i][:0]
	}
	hs.Values = hs.Values[:0]
}

func (hs *Headers) Set(key string, value string) {
	for i := 0; i < len(hs.Keys); i++ {
		if key == hs.Keys[i] {
			hs.Values[i] = hs.Values[i][:0]
			hs.Values[i] = append(hs.Values[i], value)
			return
		}
	}

	hs.Keys = append(hs.Keys, key)
	hs.Values = append(hs.Values, []string{value})
}
