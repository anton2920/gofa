/* From "textflag.h". */
#define NOSPLIT	4

TEXT ·Call0(SB), NOSPLIT, $-16
	MOVQ	fn+0(FP), AX
	CALL	AX
	MOVQ	AX, ret+8(FP)
	RET

TEXT ·Call1(SB), NOSPLIT, $-24
	MOVQ	fn+0(FP), AX
	MOVQ	a0+8(FP), DI
	CALL	AX
	MOVQ	AX, ret+16(FP)
	RET

TEXT ·Call2(SB), NOSPLIT, $-32
	MOVQ	fn+0(FP), AX
	MOVQ	a0+8(FP), DI
	MOVQ	a1+16(FP), SI
	CALL	AX
	MOVQ	AX, ret+24(FP)
	RET
