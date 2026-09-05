//go:build freebsd && amd64 && gofanostd && (gofanostd13 || gofanostd14 || gofanostd15 || gofanostd16 || gofanostdxx)
// +build freebsd,amd64,gofanostd
// +build gofanostd13 gofanostd14 gofanostd15 gofanostd16 gofanostdxx

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


// memhash_varlen(p unsafe.Pointer, h seed) uintptr
// redirects to memhash(p, h, size) using the size
// stored in the closure.
TEXT runtime·memhash_varlen(SB), 6, $32-24
	MOVQ	p+0(FP), AX
	MOVQ	h+8(FP), BX
	MOVQ	8(DX), CX
	MOVQ	AX, 0(SP)
	MOVQ	BX, 8(SP)
	MOVQ	CX, 16(SP)
	CALL	runtime·memhash(SB)
	MOVQ	24(SP), AX
	MOVQ	AX, ret+16(FP)
	RET

// memequal_varlen(a, b unsafe.Pointer) bool
TEXT runtime·memequal_varlen(SB), 6, $-17
	MOVQ	a+0(FP), SI
	MOVQ	b+8(FP), DI
	CMPQ	SI, DI
	JEQ	eq
	MOVQ	8(DX), BX    // compiler stores size at offset 8 in the closure
	LEAQ	ret+16(FP), AX
	JMP	runtime·memeqbody(SB)
eq:
	MOVB	$1, ret+16(FP)
	RET

// a in SI
// b in DI
// count in BX
// address of result byte in AX
TEXT runtime·memeqbody(SB), 6, $-0
	CMPQ	BX, $8
	JB	small

	// 64 bytes at a time using xmm registers
hugeloop:
	CMPQ	BX, $64
	JB	bigloop
	MOVOU	(SI), X0
	MOVOU	(DI), X1
	MOVOU	16(SI), X2
	MOVOU	16(DI), X3
	MOVOU	32(SI), X4
	MOVOU	32(DI), X5
	MOVOU	48(SI), X6
	MOVOU	48(DI), X7
	PCMPEQB	X1, X0
	PCMPEQB	X3, X2
	PCMPEQB	X5, X4
	PCMPEQB	X7, X6
	PAND	X2, X0
	PAND	X6, X4
	PAND	X4, X0
	PMOVMSKB X0, DX
	ADDQ	$64, SI
	ADDQ	$64, DI
	SUBQ	$64, BX
	CMPL	DX, $0xffff
	JEQ	hugeloop
	MOVB	$0, (AX)
	RET

	// 8 bytes at a time using 64-bit register
bigloop:
	CMPQ	BX, $8
	JBE	leftover
	MOVQ	(SI), CX
	MOVQ	(DI), DX
	ADDQ	$8, SI
	ADDQ	$8, DI
	SUBQ	$8, BX
	CMPQ	CX, DX
	JEQ	bigloop
	MOVB	$0, (AX)
	RET

	// remaining 0-8 bytes
leftover:
	MOVQ	-8(SI)(BX*1), CX
	MOVQ	-8(DI)(BX*1), DX
	CMPQ	CX, DX
	SETEQ	(AX)
	RET

small:
	CMPQ	BX, $0
	JEQ	equal

	LEAQ	0(BX*8), CX
	NEGQ	CX

	CMPB	SI, $0xf8
	JA	si_high

	// load at SI won't cross a page boundary.
	MOVQ	(SI), SI
	JMP	si_finish
si_high:
	// address ends in 11111xxx.  Load up to bytes we want, move to correct position.
	MOVQ	-8(SI)(BX*1), SI
	SHRQ	CX, SI
si_finish:

	// same for DI.
	CMPB	DI, $0xf8
	JA	di_high
	MOVQ	(DI), DI
	JMP	di_finish
di_high:
	MOVQ	-8(DI)(BX*1), DI
	SHRQ	CX, DI
di_finish:

	SUBQ	SI, DI
	SHLQ	CX, DI
equal:
	SETEQ	(AX)
	RET
