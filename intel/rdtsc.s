#include "textflag.h"

/* func RDTSC() intel.Cycles */
TEXT ·RDTSC(SB), NOSPLIT, $0-8
	RDTSC
	MOVL	AX, ret+0(FP)
	MOVL	DX, ret+4(FP)
	RET
