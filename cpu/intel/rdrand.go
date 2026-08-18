package intel

//go:noescape
func RDRANDW() (n uint16, ok bool)

//go:noescape
func RDRANDL() (n uint32, ok bool)

//go:noescape
func RDRANDQ() (n uint64, ok bool)
