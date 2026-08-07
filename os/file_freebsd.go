//go:build freebsd
// +build freebsd

package os

import (
	"github.com/anton2920/gofa/bits"
	"github.com/anton2920/gofa/os/posix/freebsd"
)

/* File open flags. */
const (
	OpenForReading = bits.Flags(1 << iota)
	OpenForWriting
	OpenForAppending
)

/* File creation flags. */
const (
	CreateFileIfItDoesNotExist = bits.Flags(1 << iota)
	FailCreationIfFileExists
	TruncateSizeToZero
)

func OpenOrCreateFile(path string, rw bits.Flags, creat bits.Flags, perms uint) (Handle, error) {
	var rwFlags, creatFlags uint

	/* Read/Write/Append flags. */
	switch {
	case rw.Have(OpenForReading | OpenForAppending):
		rwFlags |= freebsd.O_RDWR | freebsd.O_APPEND
	case rw.Have(OpenForAppending):
		rwFlags |= freebsd.O_WRONLY | freebsd.O_APPEND
	case rw.Have(OpenForReading | OpenForWriting):
		rwFlags |= freebsd.O_RDWR
	case rw.Have(OpenForReading):
		rwFlags |= freebsd.O_RDONLY
	case rw.Have(OpenForWriting):
		rwFlags |= freebsd.O_WRONLY
	}

	/* Create/Excl./Truncate flags. */
	if creat.Have(CreateFileIfItDoesNotExist) {
		creatFlags |= freebsd.O_CREAT
	}
	if creat.Have(FailCreationIfFileExists) {
		creatFlags |= freebsd.O_EXCL
	}
	if creat.Have(TruncateSizeToZero) {
		creatFlags |= freebsd.O_TRUNC
	}

	f, err := freebsd.Open(path, int32(rwFlags|creatFlags), uint16(perms))
	return Handle(f), err
}

//go:nosplit
func CloseHandle(f Handle) error {
	return freebsd.Close(int32(f))
}

//go:nosplit
func ReadFromFile(f Handle, buf []byte) (int, error) {
	return freebsd.Read(int32(f), buf)
}

//go:nosplit
func ReadFromFileAt(f Handle, buf []byte, offt int64) (int, error) {
	return freebsd.Pread(int32(f), buf, offt)
}

//go:nosplit
func WriteToFile(f Handle, buf []byte) (int, error) {
	return freebsd.Write(int32(f), buf)
}

//go:nosplit
func WriteToFileAt(f Handle, buf []byte, offt int64) (int, error) {
	return freebsd.Pwrite(int32(f), buf, offt)
}
