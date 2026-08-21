//go:build freebsd
// +build freebsd

package os

import (
	"github.com/anton2920/gofa/bits"
	"github.com/anton2920/gofa/context"
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

func OpenOrCreateFile(ctx *context.Context, path string, rw bits.Flags, creat bits.Flags, perms uint) (Handle, bool) {
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

	f, ok := freebsd.Open(ctx, path, int32(rwFlags|creatFlags), uint16(perms))
	return Handle(f), ok
}

//go:nosplit
func CloseHandle(ctx *context.Context, f Handle) bool {
	return freebsd.Close(ctx, int32(f))
}

//go:nosplit
func ReadFromFile(ctx *context.Context, f Handle, buf []byte) (int, bool) {
	return freebsd.Read(ctx, int32(f), buf)
}

//go:nosplit
func ReadFromFileAt(ctx *context.Context, f Handle, buf []byte, offt int64) (int, bool) {
	return freebsd.Pread(ctx, int32(f), buf, offt)
}

//go:nosplit
func WriteToFile(ctx *context.Context, f Handle, buf []byte) (int, bool) {
	return freebsd.Write(ctx, int32(f), buf)
}

//go:nosplit
func WriteToFileAt(ctx *context.Context, f Handle, buf []byte, offt int64) (int, bool) {
	return freebsd.Pwrite(ctx, int32(f), buf, offt)
}

//go:nosplit
func ResizeFile(ctx *context.Context, f Handle, size int) bool {
	return freebsd.Ftruncate(ctx, int32(f), int64(size))
}
