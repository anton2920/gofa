//go:build gofanostd
// +build gofanostd

/* NOTE(anton2920): without this function loader panics on index out of range error for some reason. */
TEXT foo<>(SB), 4, $-0
	RET

GLOBL runtime·algarray(SB), 8, $192
GLOBL runtime·writeBarrier(SB), 8, $8
GLOBL runtime·writeBarrierEnabled(SB), 24, $8
GLOBL runtime·firstmoduledata(SB), 8, $1024
GLOBL runtime·framepointer_enabled(SB), 24, $8
