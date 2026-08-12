//go:build freebsd
// +build freebsd

package os

import (
	"unsafe"

	"github.com/anton2920/gofa/bits"
	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/go/types"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/os/posix/freebsd"
	"github.com/anton2920/gofa/pointers"
)

/* TODO(anton2920): query that info on 'Init'. */
const PageSize = 4096

/* Memory protection flags. */
const (
	AllocateForNothing   = bits.Flags(freebsd.PROT_NONE)
	AllocateForReading   = bits.Flags(freebsd.PROT_READ)
	AllocateForWriting   = bits.Flags(freebsd.PROT_WRITE)
	AllocateForExecution = bits.Flags(freebsd.PROT_EXEC)
)

/* Memory allocation flags. */
const (
	FailIfRangeIsAllocated = bits.Flags(freebsd.MAP_EXCL)
	VirtualMemoryIsShared  = bits.Flags(freebsd.MAP_SHARED)
	VirtualMemoryIsPrivate = bits.Flags(freebsd.MAP_PRIVATE)
)

func AllocateVirtualMemoryStartingFrom(ctx *context.Context, base unsafe.Pointer, size int, rwx bits.Flags, extra bits.Flags) (unsafe.Pointer, bool) {
	return freebsd.Mmap(ctx, base, uint(size), int32(rwx), int32(extra)|freebsd.MAP_ANON|freebsd.MAP_FIXED, -1, 0)
}

func AllocateVirtualMemory(ctx *context.Context, size int, rwx bits.Flags) (unsafe.Pointer, bool) {
	return freebsd.Mmap(ctx, nil, uint(size), int32(rwx), freebsd.MAP_ANON|freebsd.MAP_PRIVATE, -1, 0)
}

func DeallocateVirtualMemory(ctx *context.Context, addr unsafe.Pointer, size int) bool {
	return freebsd.Munmap(ctx, addr, uint(size))
}

func CreateCircularMemoryMapping(ctx *context.Context, realSize int, virtualSize int) (unsafe.Pointer, bool) {
	/* TODO(anton2920): enable mappings backed by large pages. */
	realSize = ints.AlignUpPow2(realSize, PageSize)
	virtualSize = ints.AlignUpPow2(virtualSize, PageSize)

	/* NOTE(anton2920): first argument is SHM_ANON, cannot have that as a variable as Go's checkptr doesn't like it. */
	fd, ok := freebsd.ShmOpen2(ctx, *(*string)(unsafe.Pointer(&types.StringHeader{Data: 1, Len: 8})), freebsd.O_RDWR, 0, 0, *(*string)(unsafe.Pointer(&types.StringHeader{Data: 0, Len: 0})))
	if !ok {
		return nil, false
	}

	if !freebsd.Ftruncate(ctx, fd, int64(realSize)) {
		freebsd.Close(ctx, fd)
		return nil, false
	}

	ptr, ok := freebsd.Mmap(ctx, nil, uint(virtualSize), freebsd.PROT_NONE, freebsd.MAP_PRIVATE|freebsd.MAP_ANON, -1, 0)
	if !ok {
		freebsd.Close(ctx, fd)
		return nil, false
	}

	for i := 0; i < virtualSize; i += realSize {
		if _, ok := freebsd.Mmap(ctx, pointers.Add(ptr, uintptr(i)), uint(realSize), freebsd.PROT_READ|freebsd.PROT_WRITE, freebsd.MAP_SHARED|freebsd.MAP_FIXED, fd, 0); !ok {
			freebsd.Munmap(ctx, ptr, uint(virtualSize))
			freebsd.Close(ctx, fd)
			return nil, false
		}
	}

	freebsd.Close(ctx, fd)
	return ptr, true
}
