//go:build freebsd && 386
//+build freebsd,386

/* From "textflag.h". */
#define NOSPLIT	4


/* func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr) */
TEXT ·RawSyscall(SB),NOSPLIT, $0-28
	MOVL	trap+0(FP), AX
	LEAL		a1+4(FP), SI
	LEAL		trap+0(FP), DI
	CLD
	MOVSL
	MOVSL
	MOVSL
	INT	$0x80
	JAE	RawSyscallOK
	MOVL	$-1, r1+16(FP)
	MOVL	$-1, r2+20(FP)
	MOVL	AX, err+24(FP)
	RET
RawSyscallOK:
	MOVL	AX, r1+16(FP)
	MOVL	DX, r2+20(FP)
	MOVL	$0, err+24(FP)
	RET


/* func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) */
TEXT	·RawSyscall6(SB),NOSPLIT,$0-40
	MOVL	trap+0(FP), AX
	LEAL		a1+4(FP), SI
	LEAL		trap+0(FP), DI
	CLD
	MOVSL
	MOVSL
	MOVSL
	MOVSL
	MOVSL
	MOVSL
	INT	$0x80
	JAE	RawSyscall6OK
	MOVL	$-1, r1+28(FP)
	MOVL	$-1, r2+32(FP)
	MOVL	AX, err+36(FP)
	RET
RawSyscall6OK:
	MOVL	AX, r1+28(FP)
	MOVL	DX, r2+32(FP)
	MOVL	$0, err+36(FP)
	RET


/* func Syscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr) */
TEXT ·Syscall(SB),NOSPLIT, $0-28
	CALL	runtime·entersyscall(SB)
	MOVL	trap+0(FP), AX
	LEAL		a1+4(FP), SI
	LEAL		trap+0(FP), DI
	CLD
	MOVSL
	MOVSL
	MOVSL
	INT	$0x80
	JAE	SyscallOK
	MOVL	$-1, r1+16(FP)
	MOVL	$-1, r2+20(FP)
	MOVL	AX, err+24(FP)
	CALL	runtime·exitsyscall(SB)
	RET
SyscallOK:
	MOVL	AX, r1+16(FP)
	MOVL	DX, r2+20(FP)
	MOVL	$0, err+24(FP)
	CALL	runtime·exitsyscall(SB)
	RET


/* func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) */
TEXT	·Syscall6(SB),NOSPLIT,$0-40
	CALL	runtime·entersyscall(SB)
	MOVL	trap+0(FP), AX
	LEAL		a1+4(FP), SI
	LEAL		trap+0(FP), DI
	CLD
	MOVSL
	MOVSL
	MOVSL
	MOVSL
	MOVSL
	MOVSL
	INT	$0x80
	JAE	Syscall6OK
	MOVL	$-1, r1+28(FP)
	MOVL	$-1, r2+32(FP)
	MOVL	AX, err+36(FP)
	CALL	runtime·exitsyscall(SB)
	RET
Syscall6OK:
	MOVL	AX, r1+28(FP)
	MOVL	DX, r2+32(FP)
	MOVL	$0, err+36(FP)
	CALL	runtime·exitsyscall(SB)
	RET

