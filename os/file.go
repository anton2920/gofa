package os

import "github.com/anton2920/gofa/bits"

func OpenFile(path string, rw bits.Flags) (Handle, error) {
	return OpenOrCreateFile(path, rw, 0, 0)
}
