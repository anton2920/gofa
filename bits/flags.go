package bits

type Flags uint32

func (f *Flags) Del(g Flags) {
	*f = (*f) & (^g)
}

func (f Flags) Has(g Flags) bool {
	return (f & g) == g
}

func (f *Flags) Set(g Flags) {
	*f = (*f) | g
}

func (f *Flags) Toggle(g Flags) {
	*f = (*f) ^ g
}
