/* From "textflag.h". */
#define NOSPLIT	4

/* func CPUID(eax, ecx uint32) (reax, rebx, recx, redx uint32) */
TEXT ·CPUID(SB), NOSPLIT, $0-24
	MOVL	eax+0(FP), AX
	MOVL	ecx+4(FP), CX
	CPUID
	MOVL	AX, reax+8(FP)
	MOVL	BX, rebx+12(FP)
	MOVL	CX, recx+16(FP)
	MOVL	DX, redx+20(FP)
	RET
