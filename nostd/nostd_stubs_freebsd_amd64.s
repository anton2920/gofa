//go:build: freebsd && amd64 && gofanostd
//+build freebsd,amd64,gofanostd

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
	CALL	settls<>(SB)

	CALL	main·main(SB)

	XORL	AX, AX
	CALL	exit<>(SB)
