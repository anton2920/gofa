//go:build freebsd
// +build freebsd

package os

import (
	"unsafe"

	"github.com/anton2920/gofa/bits"
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

func AllocateVirtualMemoryStartingFrom(base unsafe.Pointer, size int, rwx bits.Flags, extra bits.Flags) (unsafe.Pointer, error) {
	return freebsd.Mmap(base, uint(size), int32(rwx), int32(extra)|freebsd.MAP_ANON|freebsd.MAP_FIXED, -1, 0)
}

func AllocateVirtualMemory(size int, rwx bits.Flags) (unsafe.Pointer, error) {
	return freebsd.Mmap(nil, uint(size), int32(rwx), freebsd.MAP_ANON|freebsd.MAP_PRIVATE, -1, 0)
}

func DeallocateVirtualMemory(addr unsafe.Pointer, size int) error {
	return freebsd.Munmap(addr, uint(size))
}

func CreateCircularMemoryMapping(realSize int, virtualSize int) (unsafe.Pointer, error) {
	/* TODO(anton2920): enable mappings backed by large pages. */
	realSize = ints.AlignUpPow2(realSize, PageSize)
	virtualSize = ints.AlignUpPow2(virtualSize, PageSize)

	/* NOTE(anton2920): first argument is SHM_ANON, cannot have that as a variable as Go's checkptr doesn't like it. */
	fd, err := freebsd.ShmOpen2(*(*string)(unsafe.Pointer(&types.StringHeader{Data: 1, Len: 8})), freebsd.O_RDWR, 0, 0, *(*string)(unsafe.Pointer(&types.StringHeader{Data: 0, Len: 0})))
	if err != nil {
		return nil, err
	}

	if err := freebsd.Ftruncate(fd, int64(realSize)); err != nil {
		freebsd.Close(fd)
		return nil, err
	}

	ptr, err := freebsd.Mmap(nil, uint(virtualSize), freebsd.PROT_NONE, freebsd.MAP_PRIVATE|freebsd.MAP_ANON, -1, 0)
	if err != nil {
		freebsd.Close(fd)
		return nil, err
	}

	for i := 0; i < virtualSize; i += realSize {
		if _, err := freebsd.Mmap(pointers.Add(ptr, uintptr(i)), uint(realSize), freebsd.PROT_READ|freebsd.PROT_WRITE, freebsd.MAP_SHARED|freebsd.MAP_FIXED, fd, 0); err != nil {
			freebsd.Munmap(ptr, uint(virtualSize))
			freebsd.Close(fd)
			return nil, err
		}
	}

	freebsd.Close(fd)
	return ptr, nil
}
