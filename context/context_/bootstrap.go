package context_

import (
	"unsafe"

	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/os"
	"github.com/anton2920/gofa/pointers"
)

var stub = make([]byte, os.PageSize)

func BootstrapWithFourSizes(ctx *context.Context, arenaSize int, fmtSize int, logSize int, errSize int) bool {
	ctx.InitWithEvenlySplitByteSlice(stub)

	arenaSize = ints.AlignUpPow2(arenaSize, int(unsafe.Alignof(int(0))))
	fmtSize = ints.AlignUpPow2(fmtSize, int(unsafe.Alignof(int(0))))
	logSize = ints.AlignUpPow2(logSize, int(unsafe.Alignof(int(0))))
	errSize = ints.AlignUpPow2(errSize, int(unsafe.Alignof(int(0))))

	size := ints.AlignUpPow2(arenaSize+fmtSize+logSize+errSize, os.PageSize)
	ptr, ok := os.AllocateVirtualMemory(ctx, size, os.AllocateForReading|os.AllocateForWriting)
	if !ok {
		ctx.NewError().S("failed to allocate requested amount of virtual memory: ").S(ctx.OldError()).S(" (code=").D(ctx.OldErrorCode()).S(")")
		return false
	}

	arenaBuf := bytes.SliceFromUnsafePointer(ptr, arenaSize)
	fmtBuf := bytes.SliceFromUnsafePointer(pointers.Add(ptr, uintptr(arenaSize)), fmtSize)
	logBuf := bytes.SliceFromUnsafePointer(pointers.Add(ptr, uintptr(arenaSize+fmtSize)), logSize)
	errBuf := bytes.SliceFromUnsafePointer(pointers.Add(ptr, uintptr(arenaSize+fmtSize+logSize)), errSize)

	ctx.InitWithFourByteSlices(arenaBuf, fmtBuf, logBuf, errBuf)
	return true
}

func BootstrapWithEvenlySplitSize(ctx *context.Context, size int) bool {
	ctx.InitWithEvenlySplitByteSlice(stub)

	size = ints.AlignUpPow2(size, os.PageSize)
	ptr, ok := os.AllocateVirtualMemory(ctx, size, os.AllocateForReading|os.AllocateForWriting)
	if !ok {
		ctx.NewError().S("failed to allocate requested amount of virtual memory: ").S(ctx.OldError()).S(" (code=").D(ctx.OldErrorCode()).S(")")
		return false
	}
	buf := bytes.SliceFromUnsafePointer(ptr, size)

	ctx.InitWithEvenlySplitByteSlice(buf)
	return true
}
