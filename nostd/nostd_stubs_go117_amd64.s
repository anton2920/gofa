//go:build: amd64 && gofanostd && gofanostd117
// +build amd64,gofanostd,gofanostd117

#include "textflag.h"
#include "nostd_stubs_intel.h"
#include "nostd_stubs_amd64.h"

// gcWriteBarrier performs a heap pointer write and informs the GC.
//
// gcWriteBarrier does NOT follow the Go ABI. It takes two arguments:
// - DI is the destination of the write
// - AX is the value being written at DI
// It clobbers FLAGS. It does not clobber any general-purpose registers,
// but may clobber others (e.g., SSE registers).
// Defined as ABIInternal since it does not use the stack-based Go ABI.
TEXT runtime·gcWriteBarrier<ABIInternal>(SB), NOSPLIT, $-0
	// Do the write.
	MOVQ	AX, (DI)
	RET

// gcWriteBarrierBX is gcWriteBarrier, but with args in DI and BX.
// Defined as ABIInternal since it does not use the stable Go ABI.
TEXT runtime·gcWriteBarrierBX<ABIInternal>(SB), NOSPLIT, $-0
	XCHGQ BX, AX
	CALL runtime·gcWriteBarrier<ABIInternal>(SB)
	XCHGQ BX, AX
	RET

// gcWriteBarrierDX is gcWriteBarrier, but with args in DI and DX.
// Defined as ABIInternal since it does not use the stable Go ABI.
TEXT runtime·gcWriteBarrierDX<ABIInternal>(SB), NOSPLIT, $-0
	XCHGQ DX, AX
	CALL runtime·gcWriteBarrier<ABIInternal>(SB)
	XCHGQ DX, AX
	RET

// gcWriteBarrierR8 is gcWriteBarrier, but with args in DI and R8.
// Defined as ABIInternal since it does not use the stable Go ABI.
TEXT runtime·gcWriteBarrierR8<ABIInternal>(SB), NOSPLIT, $-0
	XCHGQ R8, AX
	CALL runtime·gcWriteBarrier<ABIInternal>(SB)
	XCHGQ R8, AX
	RET


// memequal_varlen(a, b unsafe.Pointer) bool
TEXT runtime·memequal_varlen<ABIInternal>(SB), 4, $-17
#ifdef GOEXPERIMENT_regabiargs
	// AX = a       (want in SI)
	// BX = b       (want in DI)
	// 8(DX) = size (want in BX)
	CMPQ	AX, BX
	JNE	neq
	MOVQ	$1, AX	// return 1
	RET
neq:
	MOVQ	AX, SI
	MOVQ	BX, DI
	MOVQ	8(DX), BX    // compiler stores size at offset 8 in the closure
	JMP	memeqbody<>(SB)
#else
	MOVQ	a+0(FP), SI
	MOVQ	b+8(FP), DI
	CMPQ	SI, DI
	JEQ	eq
	MOVQ	8(DX), BX    // compiler stores size at offset 8 in the closure
	LEAQ	ret+16(FP), AX
	JMP	memeqbody<>(SB)
eq:
	MOVB	$1, ret+16(FP)
	RET
#endif

// Input:
//   a in SI
//   b in DI
//   count in BX
#ifndef GOEXPERIMENT_regabiargs
//   address of result byte in AX
#else
// Output:
//   result in AX
#endif
TEXT memeqbody<>(SB), 4, $-0
	CMPQ	BX, $8
	JB	small
	CMPQ	BX, $64
	JB	bigloop

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
#ifdef GOEXPERIMENT_regabiargs
	XORQ	AX, AX	// return 0
#else
	MOVB	$0, (AX)
#endif
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
#ifdef GOEXPERIMENT_regabiargs
	XORQ	AX, AX	// return 0
#else
	MOVB	$0, (AX)
#endif
	RET

	// remaining 0-8 bytes
leftover:
	MOVQ	-8(SI)(BX*1), CX
	MOVQ	-8(DI)(BX*1), DX
	CMPQ	CX, DX
#ifdef GOEXPERIMENT_regabiargs
	SETEQ	AX
#else
	SETEQ	(AX)
#endif
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
	// address ends in 11111xxx. Load up to bytes we want, move to correct position.
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
#ifdef GOEXPERIMENT_regabiargs
	SETEQ	AX
#else
	SETEQ	(AX)
#endif
	RET


TEXT runtime·duffzero<ABIInternal>(SB), 4, $-0
	MOVUPS	X15,(DI)
	MOVUPS	X15,16(DI)
	MOVUPS	X15,32(DI)
	MOVUPS	X15,48(DI)
	LEAQ	64(DI),DI

	MOVUPS	X15,(DI)
	MOVUPS	X15,16(DI)
	MOVUPS	X15,32(DI)
	MOVUPS	X15,48(DI)
	LEAQ	64(DI),DI

	MOVUPS	X15,(DI)
	MOVUPS	X15,16(DI)
	MOVUPS	X15,32(DI)
	MOVUPS	X15,48(DI)
	LEAQ	64(DI),DI

	MOVUPS	X15,(DI)
	MOVUPS	X15,16(DI)
	MOVUPS	X15,32(DI)
	MOVUPS	X15,48(DI)
	LEAQ	64(DI),DI

	MOVUPS	X15,(DI)
	MOVUPS	X15,16(DI)
	MOVUPS	X15,32(DI)
	MOVUPS	X15,48(DI)
	LEAQ	64(DI),DI

	MOVUPS	X15,(DI)
	MOVUPS	X15,16(DI)
	MOVUPS	X15,32(DI)
	MOVUPS	X15,48(DI)
	LEAQ	64(DI),DI

	MOVUPS	X15,(DI)
	MOVUPS	X15,16(DI)
	MOVUPS	X15,32(DI)
	MOVUPS	X15,48(DI)
	LEAQ	64(DI),DI

	MOVUPS	X15,(DI)
	MOVUPS	X15,16(DI)
	MOVUPS	X15,32(DI)
	MOVUPS	X15,48(DI)
	LEAQ	64(DI),DI

	MOVUPS	X15,(DI)
	MOVUPS	X15,16(DI)
	MOVUPS	X15,32(DI)
	MOVUPS	X15,48(DI)
	LEAQ	64(DI),DI

	MOVUPS	X15,(DI)
	MOVUPS	X15,16(DI)
	MOVUPS	X15,32(DI)
	MOVUPS	X15,48(DI)
	LEAQ	64(DI),DI

	MOVUPS	X15,(DI)
	MOVUPS	X15,16(DI)
	MOVUPS	X15,32(DI)
	MOVUPS	X15,48(DI)
	LEAQ	64(DI),DI

	MOVUPS	X15,(DI)
	MOVUPS	X15,16(DI)
	MOVUPS	X15,32(DI)
	MOVUPS	X15,48(DI)
	LEAQ	64(DI),DI

	MOVUPS	X15,(DI)
	MOVUPS	X15,16(DI)
	MOVUPS	X15,32(DI)
	MOVUPS	X15,48(DI)
	LEAQ	64(DI),DI

	MOVUPS	X15,(DI)
	MOVUPS	X15,16(DI)
	MOVUPS	X15,32(DI)
	MOVUPS	X15,48(DI)
	LEAQ	64(DI),DI

	MOVUPS	X15,(DI)
	MOVUPS	X15,16(DI)
	MOVUPS	X15,32(DI)
	MOVUPS	X15,48(DI)
	LEAQ	64(DI),DI

	MOVUPS	X15,(DI)
	MOVUPS	X15,16(DI)
	MOVUPS	X15,32(DI)
	MOVUPS	X15,48(DI)
	LEAQ	64(DI),DI

	RET

TEXT runtime·duffcopy<ABIInternal>(SB), 4, $-0
	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	RET
