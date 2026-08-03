//go:build freebsd
// +build freebsd

package os

import (
	"unsafe"

	"github.com/anton2920/gofa/bits"
	"github.com/anton2920/gofa/os/posix/freebsd"
)

/* TODO(anton2920): query that info on 'Init'. */
const PageSize = 4096

/* Memory protection flags. */
const (
	AllocateForNothing   = freebsd.PROT_NONE
	AllocateForReading   = freebsd.PROT_READ
	AllocateForWriting   = freebsd.PROT_WRITE
	AllocateForExecution = freebsd.PROT_EXEC
)

/* Memory allocation flags. */
const (
	FailIfRangeIsAllocated = bits.Flags(1 << iota)
)

func AllocateVirtualMemoryStartingFrom(base unsafe.Pointer, siz uint, rwx bits.Flags, extra bits.Flags) (unsafe.Pointer, error) {
	var allocationFlags uint

	if extra.Have(FailIfRangeIsAllocated) {
		allocationFlags |= freebsd.MAP_EXCL
	}

	return freebsd.Mmap(base, siz, int32(rwx), int32(allocationFlags|freebsd.MAP_ANON|freebsd.MAP_FIXED), -1, 0)
}

func AllocateVirtualMemory(siz uint, rwx bits.Flags) (unsafe.Pointer, error) {
	return freebsd.Mmap(nil, siz, int32(rwx), freebsd.MAP_ANON, -1, 0)
}

func DeallocateVirtualMemory(addr unsafe.Pointer, siz uint) error {
	return freebsd.Munmap(addr, siz)
}
