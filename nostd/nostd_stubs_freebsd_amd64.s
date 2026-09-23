//go:build freebsd && amd64 && gofanostd
// +build freebsd,amd64,gofanostd

#ifndef NOSPLIT
#define NOSPLIT	4
#endif // NOSPLIT

/* From 'src/runtime/sys_freebsd_amd64.s'. */
TEXT settls<>(SB), NOSPLIT, $8
	ADDQ	$16, AX	// adjust for ELF: wants to use -16(FS) and -8(FS) for g and m
	MOVQ	AX, 0(SP)
	MOVQ	SP, SI
	MOVQ	$129, DI	// AMD64_SET_FSBASE
	MOVQ	$165, AX	// sysarch
	SYSCALL
	JCC	2(PC)
	MOVL	$0xf1, 0xf1  // crash
	RET

TEXT	·exit(SB), NOSPLIT, $-0
	MOVL	code+0(FP), DI
	MOVL	$1, AX
	SYSCALL
	MOVQ	(AX), AX // crash

TEXT _rt0_amd64_freebsd(SB), NOSPLIT, $-0
	MOVQ	$·tls(SB), AX
	MOVQ	$·g(SB), BX
	MOVQ	BX, (AX)
	MOVQ	BX, 8(AX)
	CALL	settls<>(SB)

	/* TODO(anton2920): either remove the need of 'init' or make it work. */
	//CALL	main·init(SB)
	CALL	·main(SB)

	XORL	AX, AX
	CALL	·exit(SB)
