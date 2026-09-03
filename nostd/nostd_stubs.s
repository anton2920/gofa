//go:build gofanostd
//+build gofanostd

TEXT runtime·memmove(SB), 6, $-0
	JMP	·memmove(SB)

TEXT runtime·memequal64(SB), 6, $-0
	JMP	·memequal64(SB)

TEXT runtime·cmpstring(SB), 6, $-0
	JMP	·cmpstring(SB)

TEXT runtime·cmpbytes(SB), 6, $-0
	JMP	·cmpbytes(SB)

TEXT runtime·growslice(SB), 6, $-0
	JMP	runtime·gopanic(SB)

TEXT runtime·memhash(SB), 6, $-0
	JMP	·memhash(SB)

TEXT runtime·strhash(SB), 6, $-0
	JMP	·strhash(SB)

TEXT runtime·memequal(SB), 6, $-0
	JMP	·memequal(SB)

TEXT runtime·eqstring(SB), 6, $-0
	JMP	·eqstring(SB)


TEXT runtime·morestack_noctxt(SB), 6, $-0
	RET

TEXT runtime·panicslice(SB), 6, $-0
	MOVL	$66, DI
	MOVL	$1, AX
	SYSCALL

TEXT runtime·panicindex(SB), 6, $-0
	MOVL	$67, DI
	MOVL	$1, AX
	SYSCALL

TEXT runtime·panicwrap(SB), 6, $-0
	MOVL	$68, DI
	MOVL	$1, AX
	SYSCALL

TEXT runtime·writebarrierfat(SB), 6, $-0
	JMP	·writebarrierfat(SB)

TEXT runtime·writebarrierptr(SB), 6, $-0
	JMP	·writebarrierptr(SB)

TEXT runtime·writebarrierslice(SB), 6, $-0
	JMP	·writebarrierslice(SB)


GLOBL runtime·algarray(SB), 8, $192
