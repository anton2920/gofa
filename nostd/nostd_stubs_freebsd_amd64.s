//go:build freebsd && amd64 && gofanostd && (gofanostd13 || gofanostd14 || gofanostd15 || gofanostd16 || gofanostd17 || gofanostdxx)
// +build freebsd,amd64,gofanostd
// +build gofanostd13 gofanostd14 gofanostd15 gofanostd16 gofanostd17 gofanostdxx

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

TEXT runtime·gopanic(SB), 6, $-0
	MOVL	$69, DI
	MOVL	$1, AX
	SYSCALL

TEXT runtime·panic(SB), 6, $-0
	JMP	runtime·gopanic(SB)


/* From 'src/runtime/sys_freebsd_amd64.s'. */
TEXT settls<>(SB), 4, $8
	ADDQ	$16, AX	// adjust for ELF: wants to use -16(FS) and -8(FS) for g and m
	MOVQ	AX, 0(SP)
	MOVQ	SP, SI
	MOVQ	$129, DI	// AMD64_SET_FSBASE
	MOVQ	$165, AX	// sysarch
	SYSCALL
	JCC	2(PC)
	MOVL	$0xf1, 0xf1  // crash
	RET

TEXT	exit<>(SB), 4, $-0
	MOVL	AX, DI
	MOVL	$1, AX
	SYSCALL
	MOVQ	(AX), AX // crash

TEXT _rt0_amd64_freebsd(SB), 6, $-0
	MOVQ	$·tls(SB), AX
	MOVQ	$·g(SB), BX
	MOVQ	BX, (AX)
	MOVQ	BX, 8(AX)
	CALL	settls<>(SB)

	/* TODO(anton2920): either remove the need of 'init' or make it work. */
	//CALL	main·init(SB)
	CALL	main·main(SB)

	XORL	AX, AX
	CALL	exit<>(SB)
