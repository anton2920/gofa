//go:build freebsd
// +build freebsd

package os

import (
	"github.com/anton2920/gofa/bits"
	"github.com/anton2920/gofa/os/posix/freebsd"
)

const (
	CreateUnlessExists = bits.Flags(1 << iota)
	TruncateToZero
)

const (
	OpenForReading = bits.Flags(1 << iota)
	OpenForWriting
	OpenForAppending
)

func rw2PlatformFlags(flags bits.Flags) uint {
	var f uint

	switch {
	case flags.Have(OpenForReading | OpenForAppending):
		f |= freebsd.O_RDWR | freebsd.O_APPEND
	case flags.Have(OpenForAppending):
		f |= freebsd.O_WRONLY | freebsd.O_APPEND
	case flags.Have(OpenForReading | OpenForWriting):
		f |= freebsd.O_RDWR
	case flags.Have(OpenForReading):
		f |= freebsd.O_RDONLY
	case flags.Have(OpenForWriting):
		f |= freebsd.O_WRONLY
	}

	return f
}

func creat2PlatformFlags(flags bits.Flags) uint {
	var f uint

	switch {
	case flags.Have(CreateIfNotExists):
		f |= freebsd.O_CREAT | freebsd.O_EXCL
	case flags.Have(CreateAndTruncate):
		f |= freebsd.O_CREAT | freebsd.O_TRUNC
	}

	return f
}

func OpenFile(path string, rw bits.Flags) (Handle, error) {
	return OpenOrCreateFile(path, rw, 0, 0)
}

func OpenOrCreateFile(path string, rw bits.Flags, creat bits.Flags, perms uint) (Handle, error) {
	pflags := rw2PlatformFlags(rw) | creat2PlatformFlags(creat)

	f, err := freebsd.Open(path, int32(pflags), uint16(perms))
	return Handle(f), err
}

func ReadFromFile(f Handle, buf []byte) (int, error) {
	return freebsd.Read(int32(f), buf)
}
