/* From "textflag.h". */
#define NOSPLIT	4

/* func RDTSC() intel.Cycles */
TEXT ·RDTSC(SB), NOSPLIT, $0-8
	RDTSC
	MOVL	AX, ret_lo+0(FP)
	MOVL	DX, ret_hi+4(FP)
	RET
