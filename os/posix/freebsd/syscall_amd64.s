//go:build freebsd && amd64
//+build freebsd,amd64

/* From "textflag.h". */
#define NOSPLIT	4


/* func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr) */
TEXT ·RawSyscall(SB), NOSPLIT, $0-56
	MOVQ	trap+0(FP), AX
	MOVQ	a1+8(FP), DI
	MOVQ	a2+16(FP), SI
	MOVQ	a3+24(FP), DX
	SYSCALL
	JCC	RawSyscallOK
	MOVQ	$-1, r1+32(FP)
	MOVQ	$-1, r2+40(FP)
	MOVQ	AX, errno+48(FP)
	RET
RawSyscallOK:
	MOVQ	AX, r1+32(FP)
	MOVQ	DX, r2+40(FP)
	MOVQ	$0, errno+48(FP)
	RET


/* func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) */
TEXT ·RawSyscall6(SB), NOSPLIT, $0-80
	MOVQ	trap+0(FP), AX
	MOVQ	a1+8(FP), DI
	MOVQ	a2+16(FP), SI
	MOVQ	a3+24(FP), DX
	MOVQ	a4+32(FP), R10
	MOVQ	a5+40(FP), R8
	MOVQ	a6+48(FP), R9
	SYSCALL
	JCC	RawSyscall6OK
	MOVQ	$-1, r1+56(FP)
	MOVQ	$-1, r2+64(FP)
	MOVQ	AX, errno+72(FP)
	RET
RawSyscall6OK:
	MOVQ	AX, r1+56(FP)
	MOVQ	DX, r2+64(FP)
	MOVQ	$0, errno+72(FP)
	RET


/* func Syscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr) */
TEXT ·Syscall(SB), NOSPLIT, $0-56
	CALL	runtime·entersyscall(SB)
	MOVQ	trap+0(FP), AX
	MOVQ	a1+8(FP), DI
	MOVQ	a2+16(FP), SI
	MOVQ	a3+24(FP), DX
	SYSCALL
	JCC	SyscallOK
	MOVQ	$-1, r1+32(FP)
	MOVQ	$-1, r2+40(FP)
	MOVQ	AX, errno+48(FP)
	CALL	runtime·exitsyscall(SB)
	RET
SyscallOK:
	MOVQ	AX, r1+32(FP)
	MOVQ	DX, r2+40(FP)
	MOVQ	$0, errno+48(FP)
	CALL	runtime·exitsyscall(SB)
	RET


/* func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) */
TEXT ·Syscall6(SB), NOSPLIT, $0-80
	CALL	runtime·entersyscall(SB)
	MOVQ	trap+0(FP), AX
	MOVQ	a1+8(FP), DI
	MOVQ	a2+16(FP), SI
	MOVQ	a3+24(FP), DX
	MOVQ	a4+32(FP), R10
	MOVQ	a5+40(FP), R8
	MOVQ	a6+48(FP), R9
	SYSCALL
	JCC	RawSyscall6OK
	MOVQ	$-1, r1+56(FP)
	MOVQ	$-1, r2+64(FP)
	MOVQ	AX, errno+72(FP)
	CALL	runtime·exitsyscall(SB)
	RET
RawSyscall6OK:
	MOVQ	AX, r1+56(FP)
	MOVQ	DX, r2+64(FP)
	MOVQ	$0, errno+72(FP)
	CALL	runtime·exitsyscall(SB)
	RET

